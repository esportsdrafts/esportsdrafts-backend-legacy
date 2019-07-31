package internal

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	auth "github.com/barreyo/efantasy/services/auth/api"
	"github.com/barreyo/efantasy/services/auth/db"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

const defaultErrorMessage = "Authentication server error"

type AuthAPI struct {
	dbHandler *gorm.DB
}

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

func (a *AuthAPI) PerformAuth(ctx echo.Context) error {
	var newAuthClaim auth.AuthClaim

	err := ctx.Bind(&newAuthClaim)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid request format")
	}

	switch newAuthClaim.Claim {
	case "username+password":
		var account db.Account

		err := a.dbHandler.Where("username = ?", newAuthClaim.Username).First(&account).Error
		// Verify username and password
		if err != nil {
			return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid username or password")
		}

		match, err := ComparePasswordAndHash(*newAuthClaim.Password, account.Password)
		if err != nil {
			return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
		}

		if !match {
			return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid username or password")
		}

		return nil
	default:
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid authentication claim")
	}
}

func validUserNameString(name string) bool {
	characterCount := utf8.RuneCountInString(name)

	// Check max and min length
	// TODO: Make these limits globally configurable
	if characterCount < 5 || characterCount > 30 {
		return false
	}

	// Make sure only valid characters in name [a-z][0-9] and - or _
	for _, r := range name {
		if r == '_' || r == '-' {
			continue
		}
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return false
		}
	}

	// Protect test accounts
	if strings.HasPrefix(name, "test_user") {
		return false
	}
	return true
}

func validPasswordString(password string) bool {
	if len(password) < 12 || len(password) > 127 {
		return false
	}
	return true
}

// TODO: Do more validation?
func validEmailString(email string) bool {
	var re = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if len(re.FindStringIndex(email)) == 0 {
		return false
	}
	return true
}

func (a *AuthAPI) CreateAccount(ctx echo.Context) error {
	var newAccount auth.Account
	err := ctx.Bind(&newAccount)

	newUsername := strings.ToLower(newAccount.Username)

	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid request format")
	}

	if !validUserNameString(newUsername) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Username has to be between 5 and 30 characters inclusive and can only contain [a-z][0-9], underscores and dashes")
	}

	// Still in plain text at this point
	if !validPasswordString(newAccount.Password) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Password has to be between 12 and 127 characters inclusive")
	}

	if !validEmailString(newAccount.Email) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Invalid email format")
	}

	// TODO: These can be in the same query

	var count int
	// Check if username is in use
	a.dbHandler.Where(db.Account{Username: newUsername}).Count(&count)
	if count > 0 {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			fmt.Sprintf("Username '%s' already in use", newUsername))
	}

	// Check if email is in use
	a.dbHandler.Where(db.Account{Email: newAccount.Email}).Count(&count)
	if count > 0 {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			fmt.Sprintf("The provided email '%s' is already registered", newAccount.Email))
	}

	// Grab plain-text password, salt+hash it then save to DB
	hashingParams := GetDefaultHashingParams()
	hashedPassword, err := GenerateFromPassword(newAccount.Password, hashingParams)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, defaultErrorMessage)
	}

	dbAccount := &db.Account{
		Username: newUsername,
		Email:    newAccount.Email,
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
