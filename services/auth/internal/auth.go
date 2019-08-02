package internal

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	auth "github.com/barreyo/efantasy/services/auth/api"
	"github.com/barreyo/efantasy/services/auth/db"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

const defaultErrorMessage = "Authentication server error"

type AuthAPI struct {
	dbHandler *gorm.DB
}

// NewAuthAPI constructs an API client
func NewAuthAPI(dbHandler *gorm.DB) *AuthAPI {
	return &AuthAPI{
		dbHandler: dbHandler,
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

func shouldSetCookie(ctx echo.Context) bool {
	if ctx.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" {
		return true
	}
	return true
}

func writeHeaderPayloadCookie(ctx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(30 * time.Minute)
	ctx.SetCookie(cookie)
}

func writeSignatureCookie(ctx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(24 * time.Hour)
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

		if !ValidUserNameString(*newAuthClaim.Username) {
			return sendAuthAPIError(ctx, http.StatusUnprocessableEntity, "Invalid username or password")
		}

		// Still in plain text at this point
		if !ValidPasswordString(*newAuthClaim.Password) {
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

		return nil
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

	if !ValidUserNameString(newUsername) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Username has to be between 5 and 30 characters inclusive and can only contain [a-z][0-9], underscores and dashes")
	}

	// Still in plain text at this point
	if !ValidPasswordString(newPassword) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Password has to be between 12 and 127 characters inclusive")
	}

	if !ValidEmailString(newEmail) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Invalid email format")
	}

	// TODO: These queries can be in the same transaction

	// Check if username is in use
	err = a.dbHandler.Where(db.Account{Username: newUsername}).Error
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Username '%s' already in use", newUsername))
	}

	// Check if email is in use
	err = a.dbHandler.Where(db.Account{Email: newEmail}).Error
	if err != nil {
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

	// Empty JSON body with success status
	err = ctx.JSON(http.StatusCreated, map[string]int{})
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	return nil
}

// Verify takes a token and if it is associated with any user, mark the user's
// email as 'verified'.
func (a *AuthAPI) Verify(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]int{})
}
