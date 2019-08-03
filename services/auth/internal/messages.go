package internal

import (
	"encoding/json"
	"fmt"
	"time"

	efanlog "github.com/barreyo/efantasy/libs/log"
	"github.com/beanstalkd/go-beanstalk"
)

// TODO: These generic types, move to global library
type BeanstalkdClient struct {
	Address string
	Port    string
}

type BeanstalkdJob struct {
	JobType string `json:"job_type"`
}

type welcomeEmail struct {
	BeanstalkdJob
	Username         string `json:"username"`
	Email            string `json:"email"`
	VerificationCode string `json:"verification_code"`
}

const welcomeEmailJobPriority = 1
const resetEmailJobPriority = 1
const defaultJobTTR = 30 * time.Second
const defaultJobDelay = 0

// ScheduleNewUserEmail schedules a welcome email with email verification
func (bc *BeanstalkdClient) ScheduleNewUserEmail(username string, email string, verificationCode string) (uint64, error) {
	c, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%s", bc.Address, bc.Port))
	if err != nil {
		efanlog.GetLogger().Errorf("Failed to schedule welcome email for user %s", username)
		return 0, fmt.Errorf("Failed to schedule welcome email")
	}

	emailJob := welcomeEmail{
		BeanstalkdJob: BeanstalkdJob{
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
func (bc *BeanstalkdClient) SchedulePasswordResetEmail(username string, email string, resetCode string) {

}
