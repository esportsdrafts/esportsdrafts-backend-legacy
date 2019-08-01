package internal

import (
	"regexp"
	"unicode/utf8"
)

// ValidUserNameString validates a username entry according to two rules:
//   * Can only contain [a-z][0-9] and - or _
// 	 * 5 <= Length <= 30
func ValidUserNameString(name string) bool {
	characterCount := utf8.RuneCountInString(name)

	// Check max and min length
	// TODO: Make these limits globally configurable
	if characterCount < 5 || characterCount > 30 {
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

// ValidPasswordString returns if password has correct length
func ValidPasswordString(password string) bool {
	// Arbitrary upper limit. Schrugz in security.
	if len([]rune(password)) < 12 || len([]rune(password)) > 128 {
		return false
	}
	return true
}

// ValidEmailString returns true if string contains @ and a punctuation,
// more validation than that will most likely be wrong and piss off users.
func ValidEmailString(email string) bool {
	var re = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	if len(re.FindStringIndex(email)) == 0 {
		return false
	}
	return true
}
