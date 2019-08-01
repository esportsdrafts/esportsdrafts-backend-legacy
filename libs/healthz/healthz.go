package healthz

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func healthz(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func RegisterHealthChecks(router *echo.Echo) {
	router.GET("/healthz", healthz)
}
