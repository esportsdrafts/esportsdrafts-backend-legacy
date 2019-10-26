package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/matcornic/hermes/v2"
)

var /* const */ env = os.Getenv("ENV")

var /* const */ yearNow = time.Now().Year()
var /* const */ h = hermes.Hermes{
	// Optional Theme
	Theme: new(hermes.Flat),
	Product: hermes.Product{
		// Appears in header & footer of e-mails
		Name: "eFantasy",
		Link: "https://efantasy.dev/",
		// Optional product logo
		Logo:      "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		Copyright: fmt.Sprintf("Copyright © %d eFantasy. All rights reserved.", yearNow),
	},
}

// Used for local dev testing, writing the email to a local file directory
// such that our integration tests can read from it and verify email content
func writeLocalEmail(emailType, username, emailBody string) error {
	if _, err := os.Stat("/inbox"); os.IsNotExist(err) {
		os.Mkdir("/inbox", os.ModeDir)
	}
	now := time.Now()
	timeNow := now.Unix()
	fileName := fmt.Sprintf("/inbox/%d_%s_%s.html", timeNow, emailType, username)
	return ioutil.WriteFile(fileName, []byte(emailBody), 0644)
}

// SendWelcomeEmail sends an email to the user
func SendWelcomeEmail(username, userEmail, code string) error {
	email := hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				"Welcome to eFantasy!",
				"We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with eFantasy, please click here:",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Confirm your account",
						Link:  fmt.Sprintf("https://efantasy.dev/confirm?user=%s&token=%s", username, code),
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		return err
	}

	// Local dev environment cannot send emails anywhere so dump to fake inbox
	// aka a file in a folder
	if env == "development" || env == "dev" {
		return writeLocalEmail("welcome", username, emailBody)
	}

	// TODO: Call email API to actually send out the email
	emailText, err := h.GeneratePlainText(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	print(emailText)

	return nil
}

// SendResetPasswordEmail sends an email to the user
func SendResetPasswordEmail(username string, userEmail string, code string) error {
	email := hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				"You have received this email because a password reset request for your eFantasy account was recieved.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To reset your password, please click here:",
					Button: hermes.Button{
						Color: "#DC4D2F", // Optional action button color
						Text:  "Reset your password",
						Link:  fmt.Sprintf("https://efantasy.dev/reset_password?user=%s&token=%s", username, code),
					},
				},
			},
			Outros: []string{
				"If you did not make this request, please ignore this message.",
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
			Signature: "Thanks",
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		return err
	}

	// Local dev environment cannot send emails anywhere so dump to fake inbox
	// aka a file in a folder
	if env == "development" || env == "dev" {
		return writeLocalEmail("reset_password", username, emailBody)
	}

	// TODO: Call email API to actually send out the email
	emailText, err := h.GeneratePlainText(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	print(emailText)

	return nil
}
