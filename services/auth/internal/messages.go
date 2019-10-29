package internal

import (
	"encoding/json"
	"fmt"
	"time"

	beanstalkd_models "github.com/barreyo/esportsdrafts/libs/beanstalkd/models"
	efanlog "github.com/barreyo/esportsdrafts/libs/log"
	"github.com/beanstalkd/go-beanstalk"
)

// Priority 0 will be processed instantly(most urgent), higher number will be
// processed with less urgency
const (
	welcomeEmailJobPriority = 1024 * 5
	resetEmailJobPriority   = 1024 * 5
	defaultJobTTR           = 30 * time.Second
	defaultJobDelay         = 0
)

// ScheduleNewUserEmail schedules a welcome email with email verification
func ScheduleNewUserEmail(client *beanstalkd_models.Client, username string, email string, verificationCode string) (uint64, error) {
	c, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%s", client.Address, client.Port))
	if err != nil {
		efanlog.GetLogger().Errorf("Failed to schedule welcome email for user %s", username)
		return 0, fmt.Errorf("Failed to schedule welcome email")
	}

	efanlog.GetLogger().Infof("Scheduling welcome email to %s (%s) with code %s", username, email, verificationCode)

	emailJob := beanstalkd_models.WelcomeEmail{
		Job: beanstalkd_models.Job{
			JobType: "welcome_email",
		},
		Username:         username,
		Email:            email,
		VerificationCode: verificationCode,
	}

	marshalled, err := json.Marshal(emailJob)
	if err != nil {
		efanlog.GetLogger().Errorf("Failed to marshal welcome email job")
		return 0, fmt.Errorf("Failed to schedule welcome email")
	}

	id, err := c.Put(marshalled, welcomeEmailJobPriority, defaultJobDelay, defaultJobTTR)
	if err != nil {
		return 0, fmt.Errorf("Failed to schedule welcome email")
	}

	return id, nil
}

// SchedulePasswordResetEmail schedules a password reset email
func SchedulePasswordResetEmail(client *beanstalkd_models.Client, username string, email string, resetCode string) (uint64, error) {
	c, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%s", client.Address, client.Port))
	if err != nil {
		efanlog.GetLogger().Errorf("Failed to schedule password reset email for user %s", username)
		return 0, fmt.Errorf("Failed to schedule password reset email")
	}

	efanlog.GetLogger().Infof("Scheduling password reset email to %s (%s) with code %s", username, email, resetCode)

	emailJob := beanstalkd_models.ResetPasswordEmail{
		Job: beanstalkd_models.Job{
			JobType: "reset_password_email",
		},
		Username:  username,
		Email:     email,
		ResetCode: resetCode,
	}

	marshalled, err := json.Marshal(emailJob)
	if err != nil {
		efanlog.GetLogger().Errorf("Failed to marshal reset password email job")
		return 0, fmt.Errorf("Failed to schedule reset password email")
	}

	id, err := c.Put(marshalled, resetEmailJobPriority, defaultJobDelay, defaultJobTTR)
	if err != nil {
		return 0, fmt.Errorf("Failed to schedule reset password email")
	}

	return id, nil
}
