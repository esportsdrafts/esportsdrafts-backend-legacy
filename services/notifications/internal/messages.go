package internal

import (
	"encoding/json"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/barreyo/efantasy/libs/beanstalkd/models"
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
			// TODO: Fix multiple parsing of JSON
			var msg models.WelcomeEmail
			err = json.Unmarshal(body, &msg)
			if err != nil {
				logger.Warnf("Failed to parse welcome message %d, with body: %s", id, body)
				c.Release(id, ReleasePriority, ReleaseDelay)
				continue
			}
			logger.Infof("Sending welcome email")
			err = SendWelcomeEmail(msg.Username, msg.Email, msg.VerificationCode)
			if err != nil {
				logger.Warnf("Failed to send welcome email. Error: %s", err)
				c.Release(id, ReleasePriority, ReleaseDelay)
				continue
			}
			break
		case "reset_password_email":
			var msg models.ResetPasswordEmail
			err = json.Unmarshal(body, &msg)
			if err != nil {
				logger.Warnf("Failed to parse reset password message %d, with body: %s", id, body)
				c.Release(id, ReleasePriority, ReleaseDelay)
				continue
			}
			logger.Infof("Sending reset password email")
			err = SendResetPasswordEmail(msg.Username, msg.Email, msg.ResetCode)
			if err != nil {
				logger.Warnf("Failed to send password reset email. Error: %s", err)
				c.Release(id, ReleasePriority, ReleaseDelay)
				continue
			}
			break
		default:
			logger.Infof("Burying job with id %d", id)
			c.Bury(id, BuryPriority)
			continue
		}

		logger.Infof("Finished job %d. Deleting from queue...", id)

		c.Delete(id)
	}
}
