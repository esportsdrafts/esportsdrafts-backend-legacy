package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	beanstalkd "github.com/barreyo/efantasy/libs/beanstalkd"
	auth "github.com/barreyo/efantasy/services/auth/api"
	"github.com/barreyo/efantasy/services/auth/db"
	"github.com/deepmap/oapi-codegen/pkg/middleware"

	efanlog "github.com/barreyo/efantasy/libs/log"
	"github.com/barreyo/efantasy/services/auth/internal"
	"github.com/heptiolabs/healthcheck"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func registerHealthChecks(user string, password string, hostname string) {
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	connParams := fmt.Sprintf("%s:%s@tcp(%s:3306)/", user, password, hostname)
	db, _ := sql.Open("mysql", connParams)

	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <1s.
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db, 1*time.Second))
	health.AddReadinessCheck("beanstalkd", healthcheck.TCPDialCheck("beanstalkd", 1*time.Second))

	go http.ListenAndServe("0.0.0.0:8086", health)
}

func main() {
	var port = flag.Int("port", 8000, "Port to serve auth API")
	var dbHostname = flag.String("db_hostname", "mysql", "DB hostname")
	var dbUser = flag.String("db_user", "root", "DB user")
	var dbPassword = flag.String("db_password", "password", "DB password")
	var beanstalkdAddr = flag.String("beanstalkd_address", "beanstalkd", "Beanstalkd address")
	var beanstalkdPort = flag.String("beanstalkd_port", "11300", "Beanstalkd port")

	var jwtKey = flag.String("jwt_key", "", "JWT signing key, needs to be same across cluster")
	flag.Parse()

	log := efanlog.GetLogger()

	if *jwtKey == "" {
		log.Fatal("'jwt_key' missing")
	}

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

	beanstalkClient := beanstalkd.CreateBeanstalkdClient(*beanstalkdAddr, *beanstalkdPort)
	authAPI := internal.NewAuthAPI(dbHandler, beanstalkClient)

	// TODO: Attach more middlewares and move to global lib for easy use
	e := echo.New()
	e.Use(echomiddleware.RequestID())
	e.Use(middleware.OapiRequestValidator(swagger))
	e.Use(efanlog.EchoLoggingMiddleware())

	// Register routes
	auth.RegisterHandlers(e, authAPI)

	registerHealthChecks(*dbUser, *dbPassword, *dbHostname)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
