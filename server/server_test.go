package server

import "testing"

func TestStringIsCommand(t *testing.T) {
	
	if stringIsCommand("!foo") {
		// do nothing
	}
	
	if stringIsCommand("bar") {
		t.Errorf("string is not a command.")
	}
	
	if stringIsCommand("f!gate") {
		t.Errorf("string is not a command.");
	}

	if stringIsCommand("! ") {
		// do nothing
	}

	if stringIsCommand("") {
		t.Errorf("string is not a command.");
	}
}

func TestCommandString(t *testing.T) {
	
	if getCommandFromString("!help") != "help" {
		t.Errorf("Command string is incorrect.")
	}

	if getCommandFromString("!foo bars") != "foo" {
		t.Errorf("Command string is incorrect.")
	}

	if getCommandFromString("!z bars") != "z" {
		t.Errorf("Command string is incorrect.")
	}

	if getCommandFromString("!foo  !bars ") != "foo" {
		t.Errorf("Command string is incorrect.")
	}
	
	if getCommandFromString("foo") != "" {
		t.Errorf("String is not a command.")
	}
	
	if getCommandFromString("!") != "" {
		t.Errorf("Command string is incorrect.")
	}
	
}