package main

import (
	// "flag"

	efanlog "github.com/barreyo/efantasy/libs/log"
	internal "github.com/barreyo/efantasy/services/notifications/internal"
)

func main() {
	// var dbHostname = flag.String("db_hostname", "mysql", "DB hostname")
	// var dbUser = flag.String("db_user", "root", "DB user")
	// var dbPassword = flag.String("db_password", "password", "DB password")
	// var beanstalkdAddr = flag.String("beanstalkd_address", "beanstalkd", "Beanstalkd address")
	// var beanstalkdPort = flag.String("beanstalkd_port", "11300", "Beanstalkd port")
	// flag.Parse()

	log := efanlog.GetLogger()
	log.Info("Starting read loop")

	internal.RunRecieveLoop()
}
