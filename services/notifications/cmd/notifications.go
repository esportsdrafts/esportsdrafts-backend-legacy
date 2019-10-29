package main

import (
	// "flag"
	"net/http"
	"os"
	"time"

	efanlog "github.com/barreyo/esportsdrafts/libs/log"
	internal "github.com/barreyo/esportsdrafts/services/notifications/internal"
	"github.com/heptiolabs/healthcheck"
)

var /* const */ env = os.Getenv("ENV")

func registerHealthChecks() {
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck("beanstalkd", healthcheck.TCPDialCheck("beanstalkd", 1*time.Second))
	go http.ListenAndServe("0.0.0.0:8086", health)
}

func main() {
	// var dbHostname = flag.String("db_hostname", "mysql", "DB hostname")
	// var dbUser = flag.String("db_user", "root", "DB user")
	// var dbPassword = flag.String("db_password", "password", "DB password")
	// var beanstalkdAddr = flag.String("beanstalkd_address", "beanstalkd", "Beanstalkd address")
	// var beanstalkdPort = flag.String("beanstalkd_port", "11300", "Beanstalkd port")
	// flag.Parse()

	log := efanlog.GetLogger()
	log.Infof("Running in env: %s", env)

	log.Info("Registering health checks")
	registerHealthChecks()

	log.Info("Starting read loop")

	internal.RunReceiveLoop()
}
