package internal

import (
	"testing"
)

// WARNING: Contains bad words! Sensitive people be warned.
func TestIsProfane(t *testing.T) {
	words := []struct {
		input   string
		isValid bool
	}{
		{"ass", true},
		{"dick", true},
		{"nice_username", false},
		{"pelle", false},
		{"b00bs", true},
	}
	for _, table := range words {
		res := IsProfane(table.input)
		if res != table.isValid {
			t.Errorf("Checking profanity of word '%s' was incorrect, got %t, wanted %t", table.input, res, table.isValid)
		}
	}
}

func TestProfanityScore(t *testing.T) {
	if ProfanityScore("RandomWord") > 0.0 {
		t.Errorf("Remember to write tests for this function!")
	}
}
