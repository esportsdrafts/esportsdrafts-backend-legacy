package healthz

import (
	"github.com/labstack/echo/v4"
)

func healthz(ctx echo.Context) error {
	return nil
}

func RegisterHealthChecks(router echo.Echo) {
	router.GET("/healthz", healthz)
}
