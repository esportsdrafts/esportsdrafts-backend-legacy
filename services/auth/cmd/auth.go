package main

import (
	"flag"
	"fmt"
	"os"

	auth "github.com/barreyo/efantasy/services/auth/api"
	"github.com/barreyo/efantasy/services/auth/db"
	"github.com/deepmap/oapi-codegen/pkg/middleware"

	"github.com/barreyo/efantasy/libs/healthz"
	"github.com/barreyo/efantasy/services/auth/internal"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	var port = flag.Int("port", 8000, "Port to server auth API")
	var dbHostname = flag.String("db_hostname", "mysql", "DB hostname")
	var dbUser = flag.String("db_user", "root", "DB hostname")
	var dbPassword = flag.String("db_password", "password", "DB hostname")
	flag.Parse()

	swagger, err := auth.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec: %s\n", err)
		os.Exit(1)
	}
	swagger.Servers = nil

	dbHandler, err := db.CreateDBHandler(*dbHostname, *dbUser, *dbPassword, "auth_db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to DB: %s\n", err)
		os.Exit(1)
	}
	defer dbHandler.Close()

	authAPI := internal.NewAuthAPI(dbHandler)

	// TODO: Attach more middlewares and move to global lib for easy use
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.RequestID())
	e.Use(middleware.OapiRequestValidator(swagger))
	healthz.RegisterHealthChecks(*e)

	// Register routes
	auth.RegisterHandlers(e, authAPI)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
