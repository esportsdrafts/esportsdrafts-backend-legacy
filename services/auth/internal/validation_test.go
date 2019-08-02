package internal

import (
	"strings"
	"testing"
)

func TestValidUsername(t *testing.T) {
	validator := GetDefaultValidator()
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
		res := validator.ValidateUsername(table.input)
		if res != table.isValid {
			t.Errorf("Validating username '%s' was incorrect, got %t, wanted %t", table.input, res, table.isValid)
		}
	}
}

func TestValidPassword(t *testing.T) {
	validator := GetDefaultValidator()
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
		res := validator.ValidatePassword(table.input)
		if res != table.isValid {
			t.Errorf("Validating password '%s' was incorrect, got %t, wanted %t", table.input, res, table.isValid)
		}
	}
}

func TestValidEmailString(t *testing.T) {
	validator := GetDefaultValidator()
	tables := []struct {
		input   string
		isValid bool
	}{
		{"pelle", false},
		{"123312312@ds.ss", true},
		{"!@@!!@!@@!@@!!!@sd.com", false},
		{"asdsadsad+23132@gmail.com", true},
		{"a@@@@@@@@@@@@@@@@@@@a.com", false},
		{"                     abc@asv.com", false},
		{"dddd     @lds.ss", false},
		{"pelle@gmail.com", true},
		{" ", false},
		{"", false},
	}
	for _, table := range tables {
		res := validator.ValidateEmail(table.input)
		if res != table.isValid {
			t.Errorf("Validating password '%s' was incorrect, got %t, wanted %t", table.input, res, table.isValid)
		}
	}
}
