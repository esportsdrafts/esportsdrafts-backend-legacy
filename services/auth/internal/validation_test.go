package internal

import (
	"strings"
	"testing"
)

func TestValidUsernameString(t *testing.T) {
	tables := []struct {
		input   string
		isValid bool
	}{
		{"pelle", true},
		{"             _", false},
		{"________", true},
		{"a", false},
		{"username with some spaces", false},
		{"12313213213213", true},
		{strings.Repeat("a", 30), true},
		{strings.Repeat("b", 34), false},
		{"", false},
	}
	for _, table := range tables {
		res := ValidUserNameString(table.input)
		if res != table.isValid {
			t.Errorf("Validating username '%s' was incorrect, got %t, wanted %t", table.input, res, table.isValid)
		}
	}
}

func TestValidPasswordString(t *testing.T) {
	tables := []struct {
		input   string
		isValid bool
	}{
		{"pelle", false},
		{"             _", true},
		{"________", false},
		{"a", false},
		{"username with some spaces", true},
		{"12313213213213", true},
		{strings.Repeat("a", 30), true},
		{strings.Repeat("b", 129), false},
		{"", false},
	}
	for _, table := range tables {
		res := ValidPasswordString(table.input)
		if res != table.isValid {
			t.Errorf("Validating password '%s' was incorrect, got %t, wanted %t", table.input, res, table.isValid)
		}
	}
}

func TestValidEmailString(t *testing.T) {
	tables := []struct {
		input   string
		isValid bool
	}{
		{"pelle", false},
		{"123312312@d.s", true},
		{"dsadasdsa@asdasd", false},
		{"!@@!!@!@@!@@!!!@@d.com", true},
		{"asdsadsad+23132@gmail.com", true},
		{"a@@@@@@@@@@@@@@@@@@@a.com", true},
		{"                     abc@asv.com", true},
		{"dddd     @lds.s", true},
		{"", false},
	}
	for _, table := range tables {
		res := ValidPasswordString(table.input)
		if res != table.isValid {
			t.Errorf("Validating password '%s' was incorrect, got %t, wanted %t", table.input, res, table.isValid)
		}
	}
}
