package healthz

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func healthz(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

// RegisterHealthChecks adds a /healthz endpoint to an Echo context that
// always responds with 200. Basically, the bare-minimum for a health check.
// For more advanced checking you should implement a liveness probe (as well).
func RegisterHealthChecks(router *echo.Echo) {
	router.GET("/healthz", healthz)
}
