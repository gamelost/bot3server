package util

import "testing"

func TestStringIsCommand(t *testing.T) {

	if !StringIsCommand("!foo") {
		t.Error("String '!foo' was not interpreted as command.")
	}

	if StringIsCommand("bar") {
		t.Errorf("String 'bar' is not a command.")
	}

	if StringIsCommand("f!gate") {
		t.Errorf("String 'f!gate' is not a command.")
	}

	if StringIsCommand("! ") {
		t.Error("Empty command prefix should not be interpreted as a command.")
	}

	if StringIsCommand("") {
		t.Errorf("String '' is not a command.")
	}

	if StringIsCommand("9") {
		t.Errorf("String '9' is not a command.")
	}
}

func TestGetCommandFromString(t *testing.T) {

	if GetCommandFromString("!foo whatever we do here") != "foo" {
		t.Errorf("Did not return 'foo'.")
	}

	if GetCommandFromString("! what happens?") != "" {
		t.Errorf("Did not return ''.")
	}

	if GetCommandFromString("foo whatever we do here") != "" {
		t.Errorf("Should be nothing.")
	}
}
