package internal

import (
	"regexp"
	"unicode/utf8"
)

// InputValidator has functions to check email, username, and passwords for
// correct formatting. TODO: Make even more general
type InputValidator interface {
	ValidateUsername(name string) bool
	ValidateEmail(email string) bool
	ValidatePassword(password string) bool
}

// BasicValidator holds a baseline implementation for an Account input
// validator.
type BasicValidator struct {
	maxUsernameLength, minUsernameLength int
	maxPasswordLength, minPasswordLength int
}

// GetDefaultValidator creates a Validator with sane defaults.
func GetDefaultValidator() BasicValidator {
	return BasicValidator{
		minUsernameLength: 5,
		maxUsernameLength: 30,
		minPasswordLength: 12,
		maxPasswordLength: 128,
	}
}

// ValidateUsername validates a username entry according to two rules:
//   * Can only contain [a-z][0-9] and - or _
// 	 * min <= Length <= ma
func (d *BasicValidator) ValidateUsername(name string) bool {
	return validUsernameString(name, d.minUsernameLength, d.maxUsernameLength)
}

// ValidateEmail returns true if string contains @ and a punctuation,
// more validation than that will most likely be wrong and piss off users.
func (d *BasicValidator) ValidateEmail(email string) bool {
	return validEmailString(email)
}

// ValidatePassword returns if password has correct length
func (d *BasicValidator) ValidatePassword(password string) bool {
	return validPasswordString(password, d.minPasswordLength, d.maxPasswordLength)
}

// validUsernameString validates a username entry according to two rules:
//   * Can only contain [a-z][0-9] and - or _
// 	 * min <= Length <= max
func validUsernameString(name string, min int, max int) bool {
	characterCount := utf8.RuneCountInString(name)

	// Check max and min length
	// TODO: Make these limits globally configurable
	if characterCount < min || characterCount > max {
		return false
	}

	// Make sure only valid characters in name [a-z][0-9] and - or _
	for _, r := range name {
		if r == '_' || r == '-' {
			continue
		}
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return false
		}
	}

	return true
}

// validPasswordString returns if password has correct length
func validPasswordString(password string, min int, max int) bool {
	// Arbitrary upper limit. Schrugz in security.
	if len([]rune(password)) < min || len([]rune(password)) > max {
		return false
	}
	return true
}

// validEmailString returns true if string contains @ and a punctuation,
// more validation than that will most likely be wrong and piss off users.
func validEmailString(email string) bool {
	var re = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	if len(re.FindStringIndex(email)) == 0 {
		return false
	}
	return true
}
