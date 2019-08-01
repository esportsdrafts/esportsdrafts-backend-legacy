package main

import (
	"flag"
	"fmt"

	auth "github.com/barreyo/efantasy/services/auth/api"
	"github.com/barreyo/efantasy/services/auth/db"
	"github.com/deepmap/oapi-codegen/pkg/middleware"

	"github.com/barreyo/efantasy/libs/healthz"
	efanlog "github.com/barreyo/efantasy/libs/log"
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

	log := efanlog.GetLogger()

	swagger, err := auth.GetSwagger()
	if err != nil {
		log.Fatal("Error loading swagger spec: ", err)
	}
	swagger.Servers = nil

	dbHandler, err := db.CreateDBHandler(*dbHostname, *dbUser, *dbPassword, "auth_db")
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}
	defer dbHandler.Close()

	authAPI := internal.NewAuthAPI(dbHandler)

	// TODO: Attach more middlewares and move to global lib for easy use
	e := echo.New()
	e.Use(echomiddleware.RequestID())
	e.Use(middleware.OapiRequestValidator(swagger))
	e.Use(efanlog.EchoLoggingMiddleware())

	healthz.RegisterHealthChecks(e)

	// Register routes
	auth.RegisterHandlers(e, authAPI)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
