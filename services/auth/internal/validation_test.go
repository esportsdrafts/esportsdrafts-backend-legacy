package internal

import "testing"

func TestValidUsernameString(t *testing.T) {
	tables := []struct {
		input   string
		isValid bool
	}{
		{"pelle", true},
		{"             _", false},
		{"________", true},
		{"a", false},
	}
	for _, table := range tables {
		res := ValidUserNameString(table.input)
		if res != table.isValid {
			t.Errorf("Validating username '%s' was incorrect, got %t, wanted %t", table.input, res, table.isValid)
		}
	}
}
