package remindme

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestStructFromCommand1(t *testing.T) {

	r, err := ReminderStructFromCommand("", "")
	// check for proper response/error
	if !(r == nil && err == nil) {
		t.Fail()
	}
}

func TestStructFromCommand2(t *testing.T) {

	r, err := ReminderStructFromCommand("foo", "1s")
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

	r, err := ReminderStructFromCommand("test", command)
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

	r, err := ReminderStructFromCommand("test", command)
	if !(r != nil || err == nil) {
		t.Fail()
	}

	if r.Message != message {
		t.Errorf("Incorrect response message. Expected:[%s], but got:[%s]", message, r.Message)
	}
}

func TestCreateNewService(t *testing.T) {

	svc := NewRemindMeService(nil, nil)

	if svc == nil {
		t.Errorf("svc did not initialize correctly, %v", svc)
	}

	// check if data directory was created
	fiinfo, err := os.Stat("data")
	if err != nil {
		// no such file, fail
		t.Errorf("os.Stat on data dir returned error: %v", err)
	}

	if !fiinfo.IsDir() {
		t.Error("FileMode is not directory")
	}
}

func TestWriteReminder(t *testing.T) {

	svc := NewRemindMeService(nil, nil)
	if svc == nil {
		t.Errorf("svc did not initialize correctly, %v", svc)
	}

	r := &Reminder{Message: "foo", Recipient: "testuser"}
	svc.WriteReminderToDataDir(r)
}

func TestReminderOccursAt(t *testing.T) {

	currTime := time.Now()
	rem := &Reminder{CreatedOn: currTime}
	rem.Duration = 10 * time.Second

	if rem.ReminderOccursAt() != currTime.Add(10*time.Second) {
		t.Errorf("Reminder occurs at is not correct.")
	}
}

func TestReadDirectory(t *testing.T) {

	svc := NewRemindMeService(nil, nil)
	svc.LoadRemindersFromDataDirectory()
}
