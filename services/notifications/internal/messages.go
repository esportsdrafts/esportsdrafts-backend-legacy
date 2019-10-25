package internal

import (
	"time"

	"github.com/Jeffail/gabs/v2"
	efanlog "github.com/barreyo/efantasy/libs/log"
	"github.com/beanstalkd/go-beanstalk"
)

// RcvTimeout denotes time to wait for messages in seconds
const (
	RcvTimeout      = 5
	ReleasePriority = 1024 * 5
	ReleaseDelay    = 10 * time.Second
	BuryPriority    = 1024
)

func RunReceiveLoop() {
	logger := efanlog.GetLogger()

	c, err := beanstalk.Dial("tcp", "beanstalkd:11300")
	if err != nil {
		logger.Fatalf("Failed to connect to Beanstalkd, error: %s", err)
	}

	for {
		id, body, err := c.Reserve(RcvTimeout * time.Second)
		if err != nil {
			c.Release(id, ReleasePriority, ReleaseDelay)
			continue
		}

		parsed, err := gabs.ParseJSON(body)
		if err != nil {
			logger.Warnf("Failed to parse message %d, with body: %s", id, body)
			c.Release(id, ReleasePriority, ReleaseDelay)
			continue
		}

		var jobType string
		jobType, ok := parsed.Path("job_type").Data().(string)
		if !ok {
			logger.Warnf("Failed to parse message %d, with body: %s", id, body)
			c.Release(id, ReleasePriority, ReleaseDelay)
			continue
		}

		switch jobType {
		case "welcome_email":
			logger.Infof("Sending welcome email")
			break
		case "reset_password_email":
			logger.Infof("Sending reset password email")
			break
		default:
			logger.Infof("Burying job with id %d", id)
			c.Bury(id, BuryPriority)
			continue
		}

		logger.Infof("Received body: %s", body)
		logger.Infof("Fetched job %d", id)

		c.Delete(id)
	}
}
