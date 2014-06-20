package remindme

import "testing"
import "fmt"

func TestStructFromCommand1(t *testing.T) {

	r, err := ReminderStructFromCommand("")
	// check for proper response/error
	if !(r == nil && err == nil) {
		t.Fail()
	}
}

func TestStructFromCommand2(t *testing.T) {

	r, err := ReminderStructFromCommand("1s")
	if !(r == nil && err != nil) {
		t.Fail()
	}

	if err.Error() != INSUFFICENT_ARGS {
		t.Errorf("Incorrect error message. Expected:[%s], but got:[%s]", INSUFFICENT_ARGS, err.Error())
	}
}

func TestStructFromCommand3(t *testing.T) {

	duration := "1s"
	message := "foo"
	command := fmt.Sprintf("%s %s", duration, message)

	r, err := ReminderStructFromCommand(command)
	if !(r != nil && err == nil) {
		t.Fail()
	}

	if r.Message != message {
		t.Errorf("Incorrect response message. Expected:[%s], but got:[%s]", message, r.Message)
	}
}

func TestStructFromCommand4(t *testing.T) {

	duration := "1s"
	message := "alli sspasllis 1.7 fork"
	command := fmt.Sprintf("%s %s", duration, message)

	r, err := ReminderStructFromCommand(command)
	if !(r != nil || err == nil) {
		t.Fail()
	}

	if r.Message != message {
		t.Errorf("Incorrect response message. Expected:[%s], but got:[%s]", message, r.Message)
	}
}
