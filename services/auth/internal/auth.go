package internal

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	efanlog "github.com/barreyo/efantasy/libs/log"
	auth "github.com/barreyo/efantasy/services/auth/api"
	"github.com/barreyo/efantasy/services/auth/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

const defaultErrorMessage = "Authentication server error"

// TODO: Fill from env variables and/or using KMS
var jwtKey = []byte("my_secret_key")

// JWTClaims holds eFantasy auth claims. Roles array denotes what the user
// can do within the application. For example and 'admin' would have elevated
// access compared to a 'user'.
type JWTClaims struct {
	Username string   `json:"username"`
	UserID   string   `json:"user_id"`
	Roles    []string `json:"roles"`
	jwt.StandardClaims
}

// AuthAPI holds global handlers for the API like Databases.
type AuthAPI struct {
	dbHandler        *gorm.DB
	beanstalkHandler *BeanstalkdClient
	inputValidator   InputValidator
}

// NewAuthAPI constructs an API client
func NewAuthAPI(dbHandler *gorm.DB, bClient *BeanstalkdClient) *AuthAPI {
	return &AuthAPI{
		dbHandler:        dbHandler,
		beanstalkHandler: bClient,
		inputValidator: &BasicValidator{
			maxUsernameLength: 30,
			minUsernameLength: 5,
			maxPasswordLength: 128,
			minPasswordLength: 12,
		},
	}
}

func sendAuthAPIError(ctx echo.Context, code int, message string) error {
	authError := auth.Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, authError)
	return err
}

// TODO: Move this into a shared auth lib
func hasRequestedWithHeader(ctx echo.Context) bool {
	if ctx.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
		return true
	}
	return false
}

func getAuthTokenFromHeader(ctx echo.Context) string {
	return ctx.Request().Header.Get("Authentication")
}

func writeHeaderPayloadCookie(ctx echo.Context, header string) {
	cookie := new(http.Cookie)
	cookie.Name = "header.payload"
	cookie.Value = header

	// Protect from sending over HTTP
	cookie.Secure = true

	// TODO: Make globally configurable
	cookie.Expires = time.Now().Add(60 * time.Minute)
	ctx.SetCookie(cookie)
}

func writeSignatureCookie(ctx echo.Context, signature string) {
	cookie := new(http.Cookie)
	cookie.Name = "signature"
	cookie.Value = signature

	// Protect from sending over HTTP
	cookie.Secure = true

	// Block JS from reading this cookie
	cookie.HttpOnly = true
	ctx.SetCookie(cookie)
}

// PerformAuth performs an authentication request the auth path is based on
// what claim the caller makes.
func (a *AuthAPI) PerformAuth(ctx echo.Context) error {
	var newAuthClaim auth.AuthClaim

	err := ctx.Bind(&newAuthClaim)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid request format")
	}

	switch newAuthClaim.Claim {
	case "username+password":
		var account db.Account

		if !a.inputValidator.ValidateUsername(*newAuthClaim.Username) {
			return sendAuthAPIError(ctx, http.StatusUnprocessableEntity, "Invalid username or password")
		}

		// Still in plain text at this point
		if !a.inputValidator.ValidatePassword(*newAuthClaim.Password) {
			return sendAuthAPIError(ctx, http.StatusUnprocessableEntity, "Invalid username or password")
		}

		err := a.dbHandler.Where("username = ?", newAuthClaim.Username).First(&account).Error
		// Verify username and password
		if err != nil {
			return sendAuthAPIError(ctx, http.StatusUnauthorized, "Invalid username or password")
		}

		match, err := ComparePasswordAndHash(*newAuthClaim.Password, account.Password)
		if err != nil {
			return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
		}

		if !match {
			return sendAuthAPIError(ctx, http.StatusUnauthorized, "Invalid username or password")
		}

		issuedTime := time.Now()
		expirationTime := issuedTime.Add(60 * time.Minute)

		// Create the JWT claims, which includes the username and expiry time
		claims := &JWTClaims{
			Username: account.Username,
			UserID:   account.ID.String(),
			// Would be something more useful depending on the user type
			Roles: []string{"user"},
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
				// Can be used to blacklist in the future. Needs to hold state
				// in that case :/
				Id:       uuid.NewV4().String(),
				IssuedAt: issuedTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
		}

		// Web client so set cookies instead of returning
		if hasRequestedWithHeader(ctx) {
			parts := strings.Split(tokenString, ".")
			signature := parts[2]
			headerPayload := strings.Join(parts[0:2], ".")
			writeSignatureCookie(ctx, signature)
			writeHeaderPayloadCookie(ctx, headerPayload)
			return ctx.JSON(http.StatusOK, map[string]int{})
		}

		result := auth.JWT{
			AccessToken: tokenString,
			ExpiresIn:   int(expirationTime.Unix()),
		}
		// Otherwise just give token
		return ctx.JSON(http.StatusOK, result)
	default:
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid authentication claim")
	}
}

// CreateAccount creates a new Account
func (a *AuthAPI) CreateAccount(ctx echo.Context) error {
	var newAccount auth.Account
	err := ctx.Bind(&newAccount)

	newUsername := strings.ToLower(newAccount.Username)
	newEmail := newAccount.Email
	newPassword := newAccount.Password

	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid request format")
	}

	if !a.inputValidator.ValidateUsername(newUsername) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Username has to be between 5 and 30 characters inclusive and can only contain [a-z][0-9], underscores and dashes")
	}

	// Still in plain text at this point
	if !a.inputValidator.ValidatePassword(newPassword) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Password has to be between 12 and 127 characters inclusive")
	}

	if !a.inputValidator.ValidateEmail(newEmail) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Invalid email format")
	}

	// TODO: These queries can be in the same transaction
	var count int
	var usernameCheck db.Account

	// Check if username is in use
	// Count vs first? Maybe simplify by not creating a struct, string
	// query instead.
	a.dbHandler.Where("username = ?", newUsername).First(&usernameCheck).Count(&count)
	if count != 0 {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Username '%s' already in use", newUsername))
	}

	var emailCheck db.Account
	// Check if email is in use
	a.dbHandler.Where("email = ?", newEmail).First(&emailCheck).Count(&count)
	if count != 0 {
		// Information leak, someone could spam and figure out which emails
		// are registered in the system.
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			fmt.Sprintf("The provided email is already registered"))
	}

	// Grab plain-text password, salt+hash it then save to DB
	hashingParams := GetDefaultHashingParams()
	hashedPassword, err := GenerateFromPassword(newPassword, hashingParams)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	// TODO: Encrypt all fields except UserID
	dbAccount := &db.Account{
		Username: newUsername,
		Email:    newEmail,
		Password: hashedPassword,
	}

	err = a.dbHandler.Save(&dbAccount).Error
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	expirationTime := time.Now().Add(48 * time.Hour)
	verifyCode := &db.EmailVerificationCode{
		UserID:    dbAccount.ID.String(),
		ExpiresAt: expirationTime,
	}

	// TODO: Transaction with creating User
	err = a.dbHandler.Save(&verifyCode).Error
	if err != nil {
		efanlog.GetLogger().Infof("Failed to create email verification token")
	} else {
		go a.beanstalkHandler.ScheduleNewUserEmail(dbAccount.Username, dbAccount.Email, verifyCode.ID.String())
	}

	// Empty JSON body with success status
	return ctx.JSON(http.StatusCreated, map[string]int{})
}

// Verify takes a token and if it is associated with any user, mark the user's
// email as 'verified'.
func (a *AuthAPI) Verify(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]int{})
}
