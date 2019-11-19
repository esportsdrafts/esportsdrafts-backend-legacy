package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	beanstalkd "github.com/esportsdrafts/esportsdrafts/libs/beanstalkd"
	auth "github.com/esportsdrafts/esportsdrafts/services/auth/api"
	"github.com/esportsdrafts/esportsdrafts/services/auth/db"

	efanlog "github.com/esportsdrafts/esportsdrafts/libs/log"
	"github.com/esportsdrafts/esportsdrafts/services/auth/internal"
	"github.com/heptiolabs/healthcheck"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func registerHealthChecks(user string, password string, hostname string) {
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	connParams := fmt.Sprintf("%s:%s@tcp(%s:3306)/", user, password, hostname)
	db, _ := sql.Open("mysql", connParams)

	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db, 2*time.Second))
	// health.AddReadinessCheck("beanstalkd", healthcheck.TCPDialCheck("beanstalkd", 2*time.Second))

	go http.ListenAndServe("0.0.0.0:8086", health)
}

func main() {
	var port = flag.Int("port", 8000, "Port to serve auth API")
	var dbHostname = flag.String("db_hostname", "mysql", "DB hostname")
	var dbUser = flag.String("db_user", "root", "DB user")
	var dbPassword = flag.String("db_password", "password", "DB password")
	var beanstalkdAddr = flag.String("beanstalkd_address", "beanstalkd", "Beanstalkd address")
	var beanstalkdPort = flag.String("beanstalkd_port", "11300", "Beanstalkd port")
	flag.Parse()

	jwtKey := os.Getenv("JWT_KEY")
	log := efanlog.GetLogger()

	if jwtKey == "" {
		log.Fatal("'JWT_KEY' not found in environment")
	}

	swagger, err := auth.GetSwagger()
	if err != nil {
		log.Fatal("Error loading swagger spec: ", err)
	}
	swagger.Servers = nil

	log.Info("Connecting to DB...")
	dbHandler, err := db.CreateDBHandler(*dbHostname, *dbUser, *dbPassword, "auth_db")
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}
	defer dbHandler.Close()

	beanstalkClient := beanstalkd.CreateBeanstalkdClient(*beanstalkdAddr, *beanstalkdPort)
	authAPI := internal.NewAuthAPI(dbHandler, beanstalkClient, []byte(jwtKey))

	// TODO: Attach more middlewares and move to global lib for easy use
	e := echo.New()
	e.Use(echomiddleware.RequestID())
	e.Use(middleware.OapiRequestValidator(swagger))
	e.Use(efanlog.EchoLoggingMiddleware())

	// Register routes
	auth.RegisterHandlers(e, authAPI)

	log.Info("Registering health checks...")
	registerHealthChecks(*dbUser, *dbPassword, *dbHostname)

	log.Infof("Initialization done. Serving on port %d...", *port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
