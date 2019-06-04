package internal

import (
	"github.com/jinzhu/gorm"
	"github.com/efantasy/auth/api"
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

func (a *AuthAPI)AuthUser(ctx echo.Context) error {
	return nil
}

func (a *AuthAPI)RefreshToken(ctx echo.Context) error {
	return nil
}

func (a *AuthAPI)CreateAccount(ctx echo.Context) error {
	return nil
}
