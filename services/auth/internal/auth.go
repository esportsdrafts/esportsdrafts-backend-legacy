package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	auth "github.com/efantasy/auth/api"
	"github.com/efantasy/auth/db"
	"github.com/labstack/echo/v4"
)

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

func (a *AuthAPI)PerformAuth(ctx echo.Context) error {
	return nil
}

func validUserNameString(name string) bool {
	return true
}

func validPasswordString(password string) bool {
	return true
}

func validEmailString(email string) bool {
	return true
}

func (a *AuthAPI)CreateAccount(ctx echo.Context) error {
	var newAccount auth.Account
	err := ctx.Bind(&newAccount)

	newUsername := strings.ToLower(newAccount.Username)

	if err != nil {
		return sendAuthAPIError(ctx, http.StatusBadRequest, "Invalid request format")
	}

	if !validUserNameString(newUsername) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Username has contain 5 or more characters and can only contain [a-z][A-Z][0-9]")
	}

	// Still in plain text at this point
	if !validPasswordString(newAccount.Password) {
		return sendAuthAPIError(ctx, http.StatusBadRequest,
			"Password has to be between 12 and 160 characters inclusive")
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

	// Grab plain-text password and hash it then save to DB
	hashingParams := GetDefaultHashingParams()
	hashedPassword, err := GenerateFromPassword(newAccount.Password, hashingParams)
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, "Failed to create new Account")
	}

	dbAccount := &db.Account {
		Username: newUsername,
		Email: newAccount.Email,
		Password: hashedPassword,
	}

	err = a.dbHandler.Save(&dbAccount).Error
	if err != nil {
		return sendAuthAPIError(ctx, http.StatusInternalServerError, "Failed to create new Account")
	}

	err = ctx.JSON(http.StatusCreated, auth.JWT{})
	if err != nil {
		return err
	}

	return nil
}
