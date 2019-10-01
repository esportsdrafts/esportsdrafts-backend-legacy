package internal

import (
	"time"

	efanlog "github.com/barreyo/efantasy/libs/log"
	"github.com/beanstalkd/go-beanstalk"
)

// RcvTimeout denotes time to wait for messages in seconds
const RcvTimeout = 5

func RunRecieveLoop() {
	logger := efanlog.GetLogger()

	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		logger.Fatalf("Failed to connect to Beanstalkd, error: %s", err)
	}

	for {
		id, body, err := c.Reserve(RcvTimeout * time.Second)
		if err != nil {
			if _, ok := err.(beanstalk.ConnError); !ok {
				continue
			}
			return
		}

		logger.Infof("Received body: %s", body)
		logger.Infof("Fetched job %d", id)
	}
}
