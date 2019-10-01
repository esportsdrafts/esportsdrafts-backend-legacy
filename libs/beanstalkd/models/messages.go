package models

type WelcomeEmail struct {
	Job
	Username         string `json:"username"`
	Email            string `json:"email"`
	VerificationCode string `json:"verification_code"`
}

type ResetPasswordEmail struct {
	Job
	Username  string `json:"username"`
	Email     string `json:"email"`
	ResetCode string `json:"reset_code"`
}
