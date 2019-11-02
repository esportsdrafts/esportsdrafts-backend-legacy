package internal

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	beanstalkd_models "github.com/barreyo/esportsdrafts/libs/beanstalkd/models"
	efanlog "github.com/barreyo/esportsdrafts/libs/log"
	auth "github.com/barreyo/esportsdrafts/services/auth/api"
	"github.com/barreyo/esportsdrafts/services/auth/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultErrorMessage = "Authentication server error"
	maxUsernameLength   = 30
	minUsernameLength   = 5
	maxPasswordLength   = 128
	minPasswordLength   = 12
)

// JWTClaims holds esportsdrafts auth claims. Roles array denotes what the user
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
	beanstalkHandler *beanstalkd_models.Client
	inputValidator   InputValidator
	jwtKey           []byte
}

// NewAuthAPI constructs an API client
func NewAuthAPI(dbHandler *gorm.DB, bClient *beanstalkd_models.Client, jwtKey []byte) *AuthAPI {
	return &AuthAPI{
		dbHandler:        dbHandler,
		beanstalkHandler: bClient,
		inputValidator: &BasicValidator{
			maxUsernameLength: maxUsernameLength,
			minUsernameLength: minUsernameLength,
			maxPasswordLength: maxPasswordLength,
			minPasswordLength: minPasswordLength,
		},
		jwtKey: jwtKey,
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
	return ctx.Request().Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func getAuthTokenFromHeader(ctx echo.Context) (string, error) {
	headerContent := ctx.Request().Header.Get("Authorization")
	headerContent = strings.TrimSpace(headerContent)
	if strings.HasPrefix(headerContent, "Bearer") {
		runes := []rune(headerContent)
		if len(runes) <= 7 {
			return "", fmt.Errorf("Auth header not found")
		}
		return strings.TrimSpace(string(runes[6:])), nil
	}
	return "", fmt.Errorf("Auth header not found")
}

func writeHeaderPayloadCookie(ctx echo.Context, header string, expiry time.Duration) {
	cookie := new(http.Cookie)
	cookie.Name = "header.payload"
	cookie.Value = header

	// Protect from sending over HTTP
	cookie.Secure = true
	cookie.Expires = time.Now().Add(expiry)
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

// GenerateAuthToken generates a auth token with provided claims
func GenerateAuthToken(claims *JWTClaims, expiry time.Duration, jwtKey []byte) (string, time.Time, error) {
	issuedTime := time.Now()
	expirationTime := issuedTime.Add(expiry)
	claims.StandardClaims = jwt.StandardClaims{
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: expirationTime.Unix(),
		// Can be used to blacklist in the future. Needs to hold state
		// in that case :/
		Id:       uuid.NewV4().String(),
		IssuedAt: issuedTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	res, err := token.SignedString(jwtKey)
	return res, expirationTime, err
}

// PerformAuth performs an authentication request the auth path is based on
// what claim the caller makes.
func (a *AuthAPI) PerformAuth(ctx echo.Context) error {
	var newAuthClaim auth.AuthClaim

	err := ctx.Bind(&newAuthClaim)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid request format")
	}
	logger := efanlog.GetLogger()

	switch newAuthClaim.Claim {
	case "username+password":
		var account db.Account

		err := a.dbHandler.Where("username = ?", newAuthClaim.Username).First(&account).Error
		// Verify username and password
		if err != nil {
			logger.Info("Username not found")
			return sendAuthAPIError(ctx, http.StatusUnauthorized, "Invalid username or password")
		}

		match, err := ComparePasswordAndHash(*newAuthClaim.Password, account.Password)
		if err != nil {
			logger.Info("Error hashing and comparing")
			return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
		}

		if !match {
			logger.Info("Username and password did not match")
			return sendAuthAPIError(ctx, http.StatusUnauthorized, "Invalid username or password")
		}

		var roles []string
		if !account.IsEmailVerified() {
			roles = append(roles, "email_verify")
		} else {
			roles = append(roles, "user")
		}

		// Create the JWT claims, which includes the username and expiry time
		claims := &JWTClaims{
			Username: account.Username,
			UserID:   account.ID.String(),
			// Would be something more useful depending on the user type
			Roles: roles,
		}

		tokenString, expirationTime, err := GenerateAuthToken(claims, 60*time.Minute, a.jwtKey)
		if err != nil {
			logger.Info(err)
			// If there is an error in creating the JWT return an internal server error
			return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
		}

		result := auth.JWT{}
		result.AccessToken = tokenString
		result.ExpiresIn = int(expirationTime.Unix())

		// Web client so set cookies
		if hasRequestedWithHeader(ctx) {
			// This can be done more efficiently by knowing exact indices of dots
			// Doing exact locations is more accurate and correct. A token should
			// always be perfectly validated
			parts := strings.Split(tokenString, ".")
			if len(parts) != 3 {
				return sendAuthAPIError(ctx, http.StatusInternalServerError, "Failed to generate auth token")
			}
			signature := parts[2]
			headerPayload := strings.Join(parts[0:2], ".")
			writeSignatureCookie(ctx, signature)

			// TODO: Configure cookie timeout globally
			writeHeaderPayloadCookie(ctx, headerPayload, 60*time.Minute)
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
	var usernameCheck db.Account
	// Check if username is in use
	err = a.dbHandler.Where("username = ?", newUsername).First(&usernameCheck).Error
	if err == nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Username '%s' already in use", newUsername))
	}

	// Important to add new reference, otherwise the query will check
	// if email + username match instead of only email
	var emailCheck db.Account
	// Check if email is in use
	err = a.dbHandler.Where("email = ?", newEmail).First(&emailCheck).Error
	if err == nil {
		// Information leak, someone could spam and figure out which emails
		// are registered in the system.
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			fmt.Sprint("The provided email is already registered"))
	}

	// Grab plain-text password, salt+hash it then save to DB
	hashingParams := GetDefaultHashingParams()
	hashedPassword, err := GenerateFromPassword(newPassword, hashingParams)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	currentTime := time.Now()
	// TODO: Encrypt all fields except UserID
	dbAccount := &db.Account{
		Username:        newUsername,
		Email:           newEmail,
		Password:        hashedPassword,
		AcceptedTermsAt: &currentTime,
	}

	err = a.dbHandler.Save(&dbAccount).Error
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	expirationTime := time.Now().Add(48 * time.Hour)
	verifyCode := &db.EmailVerificationCode{
		UserID:    dbAccount.ID,
		ExpiresAt: expirationTime,
	}

	err = a.dbHandler.Save(&verifyCode).Error
	if err != nil {
		efanlog.GetLogger().Info("Failed to create email verification token")
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}
	go ScheduleNewUserEmail(a.beanstalkHandler, dbAccount.Username, dbAccount.Email, verifyCode.ID.String())

	// Empty JSON body with success status
	return ctx.JSON(http.StatusCreated, map[string]int{})
}

// Verify takes a token and if it is associated with any user, mark the user's
// email as 'verified'.
func (a *AuthAPI) Verify(ctx echo.Context) error {
	var request auth.EmailVerification
	err := ctx.Bind(&request)

	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid request format")
	}

	var account db.Account
	err = a.dbHandler.Where("username = ?", request.Username).First(&account).Error
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "User not found")
	}

	if account.IsEmailVerified() {
		return ctx.JSON(http.StatusOK, map[string]int{})
	}

	var token db.EmailVerificationCode
	err = a.dbHandler.Where("id = ? AND user_id = ?", request.Token, account.ID).First(&token).Error
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Token not found")
	}

	if token.ExpiresAt.Before(time.Now()) {
		a.dbHandler.Delete(&token)
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Token has expired")
	}

	// Set account as verified and delete all tokens
	err = account.VerifyEmail(a.dbHandler)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	return ctx.JSON(http.StatusOK, map[string]int{})
}

// Check takes a username and verifies if it is already registered or not.
// Useful endpoint for frontend to do validation in registration form.
func (a *AuthAPI) Check(ctx echo.Context, params auth.CheckParams) error {
	if params.Username != nil {
		if !a.inputValidator.ValidateUsername(*params.Username) {
			return ctx.JSON(http.StatusUnauthorized, map[string]int{})
		}
		var usernameCheck db.Account
		err := a.dbHandler.Where("username = ?", params.Username).First(&usernameCheck).Error
		if err != nil {
			return ctx.JSON(http.StatusOK, map[string]int{})
		}
	}
	return ctx.JSON(http.StatusUnauthorized, map[string]int{})
}

// Passwordresetrequest initiates a password reset
func (a *AuthAPI) Passwordresetrequest(ctx echo.Context) error {
	var request auth.PasswordResetRequest
	err := ctx.Bind(&request)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid request format")
	}

	var account db.Account
	err = a.dbHandler.Where("username = ? AND email = ?", request.Username, request.Email).First(&account).Error

	// Always give a 200, we do not wanna reveal if the email is registered or not
	// with this username
	if err != nil {
		return ctx.JSON(http.StatusOK, map[string]int{})
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	verifyCode := &db.PasswordResetToken{
		UserID:    account.ID,
		ExpiresAt: expirationTime,
	}

	err = a.dbHandler.Save(&verifyCode).Error
	if err != nil {
		efanlog.GetLogger().Info("Failed to create password reset token")
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	go SchedulePasswordResetEmail(a.beanstalkHandler, account.Username, account.Email, verifyCode.ID.String())

	return ctx.JSON(http.StatusOK, map[string]int{})
}

// Passwordresetverify takes username, token and a new password. If the token
// matches with the password reset request the password for the account is
// changed to the supplied one in the request.
func (a *AuthAPI) Passwordresetverify(ctx echo.Context) error {
	var request auth.PasswordResetVerify
	err := ctx.Bind(&request)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid request format")
	}

	// TODO: Pull error message from validator to make sure formatting is consistent
	if !a.inputValidator.ValidatePassword(request.Password) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Password has to be between 12 and 127 characters inclusive")
	}

	var account db.Account
	// TODO: add email as well?
	err = a.dbHandler.Where("username = ?", request.Username).First(&account).Error
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, fmt.Sprint("Username not found"))
	}

	var token db.PasswordResetToken
	err = a.dbHandler.Where("id = ? AND user_id = ?", request.Token, account.ID).First(&token).Error
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid token provided")
	}

	// TODO: Remove expiresAt field and just use creation time to diff with
	// some value
	if token.ExpiresAt.Before(time.Now()) {
		a.dbHandler.Delete(&token)
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Token has expired")
	}

	hashingParams := GetDefaultHashingParams()
	hashedPassword, err := GenerateFromPassword(request.Password, hashingParams)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	err = a.dbHandler.Model(&account).Update("password_hash", hashedPassword).Error
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	return ctx.JSON(http.StatusOK, map[string]int{})
}
