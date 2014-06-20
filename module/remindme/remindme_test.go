package remindme

import "testing"

func TestStructFromCommand1(t *testing.T) {

	r, err := ReminderStructFromCommand("")
	if r != nil || err != nil {
		t.Fail()
	}
}

func TestStructFromCommand2(t *testing.T) {

	r, err := ReminderStructFromCommand("1s")
	if r != nil || err == nil {
		t.Fail()
	}
}

func TestStructFromCommand3(t *testing.T) {

	r, err := ReminderStructFromCommand("1s foo")
	if r == nil || err != nil {
		t.Fail()
	}
}

func TestStructFromCommand4(t *testing.T) {

	r, err := ReminderStructFromCommand(".5m foo")
	if r == nil || err != nil {
		t.Fail()
	}
}
