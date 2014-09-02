package remindme

import (
	iniconf "code.google.com/p/goconf/conf"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"github.com/twinj/uuid"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode"
)

// set a max duration
const MAXDURATION = time.Hour * 24 * 7

// set a min duration
const MINDURATION = time.Second * 2

const DATADIR_NAME = "data"

// string messages
const INSUFFICENT_ARGS = "Insufficent number of arguments provided.  Need to provide a duration and message."

func NewRemindMeService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) *RemindMeService {
	newSvc := &RemindMeService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	//newSvc.createDataDirectory()
	newSvc.Reminders = &ReminderList{}
	newSvc.Reminders.ReminderMap = make(map[string]*Reminder)
	//newSvc.loadRemindersFromDataDirectory()

	// switch format of uuid
	uuid.SwitchFormat(uuid.Clean)

	return newSvc
}

type RemindMeService struct {
	server.BotHandlerService
	Reminders *ReminderList
}

func (svc *RemindMeService) createDataDirectory() {

	// check if data directory was created
	fiinfo, err := os.Stat(DATADIR_NAME)
	if err == nil && fiinfo.IsDir() {
		// no such file, fail
		return
	} else {
		// create data directory
		err := os.Mkdir(DATADIR_NAME, 0777)
		if err != nil {
			panic("Unable to create data directory.")
		}
	}
}

func (svc *RemindMeService) loadRemindersFromDataDirectory() {

	// clear out the map
	svc.Reminders.ReminderMap = make(map[string]*Reminder)

	files, err := ioutil.ReadDir(DATADIR_NAME)
	if err != nil {
		panic(err)
	}

	// iterate through all files
	for _, f := range files {
		// ignore files that dont end with .json
		if strings.HasSuffix(f.Name(), ".json") {
			rem := svc.FileToReminder(f.Name())
			svc.Reminders.AddReminder(rem)
		}
	}
}

func (svc *RemindMeService) FileToReminder(filename string) *Reminder {

	rem := &Reminder{}
	bytes, err := ioutil.ReadFile(DATADIR_NAME + "/" + filename)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, rem)

	fmt.Printf("Reminder: %+v\n", rem)
	return rem
}

func (svc *RemindMeService) WriteReminderToDataDir(rem *Reminder) {

	data, err := json.Marshal(rem)
	if err != nil {
		panic(err.Error())
	}
	filename := fmt.Sprintf("%s.json", rem.ReminderIdentity())
	ioutil.WriteFile(DATADIR_NAME+"/"+filename, data, 0644)
}

func (svc *RemindMeService) RemoveReminderFromDataDir(rem *Reminder) error {

	filename := fmt.Sprintf("%s.json", rem.ReminderIdentity())
	err := os.Remove(DATADIR_NAME + "/" + filename)
	return err
}

func (svc *RemindMeService) DispatchRequest(botRequest *server.BotRequest) {

	cmd := botRequest.LineTextWithoutCommand()
	botResponse := svc.CreateBotResponse(botRequest)
	rem, err := reminderStructFromCommand(botRequest.Nick, cmd)

	if err != nil {
		botResponse.SetSingleLineResponse(fmt.Sprintf("Bloop. Your request could not be parsed: %s", err.Error()))
	} else {

		// nil reminder triggers status update instead
		if rem == nil {

			reminderCount := len(svc.Reminders.ReminderMap)

			if reminderCount == 0 {
				botResponse.SetSingleLineResponse("No reminders in queue right now.")
			} else {
				responses := make([]string, reminderCount+1)
				responses[0] = "Reminders in queue..."
				var counter = 1
				for _, r := range svc.Reminders.ReminderMap {

					responses[counter] = fmt.Sprintf("[%d]: %s, reminder will occur on: %v, (%v from now)", counter, r.Message, r.ReminderOccursAt(), r.DurationUntilReminder())
					counter++
				}
				botResponse.SetMultipleLineResponse(responses)
			}

		} else if rem.Duration < 0 {
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s, only your mom would ask you to do something in the past. You're lame.", botRequest.Nick))
		} else if rem.Duration < MINDURATION {
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s, I dont work that fast!", botRequest.Nick))
		} else if rem.Duration > MAXDURATION {
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s, really? Maybe you should use a calendar instead.  Durations less than a week please.", botRequest.Nick))
		} else {
			err := svc.ScheduleReminder(rem, botRequest)
			if err == nil {
				botResponse.SetSingleLineResponse("I'll remind ya, m8!")
			}
		}
	}

	svc.PublishBotResponse(botResponse)
}

func (svc *RemindMeService) ScheduleReminder(rem *Reminder, botRequest *server.BotRequest) error {

	svc.Reminders.AddReminder(rem)
	//svc.WriteReminderToDataDir(rem)
	// set the afterfunc
	time.AfterFunc(rem.Duration, func() {
		botResponse := svc.CreateBotResponse(botRequest)
		botResponse.SetSingleLineResponse(fmt.Sprintf("%s, you asked me to remind you: %s", rem.Recipient, rem.Message))
		svc.PublishBotResponse(botResponse)
		svc.RemoveReminder(rem)
	})

	return nil
}

func (svc *RemindMeService) RemoveReminder(rem *Reminder) {
	svc.Reminders.RemoveReminder(rem)
	//svc.RemoveReminderFromDataDir(rem)
}

type Reminder struct {
	Duration  time.Duration
	CreatedOn time.Time
	Message   string
	Recipient string
	Identity  string
}

func (rem *Reminder) ReminderIdentity() string {

	if rem.Identity == "" {
		uuid.SwitchFormat(uuid.Clean)
		rem.Identity = fmt.Sprintf("reminder-%s-%s", rem.Recipient, uuid.NewV1().String())
	}
	return rem.Identity
}

func (rem *Reminder) ReminderOccursAt() time.Time {
	return rem.CreatedOn.Add(rem.Duration)
}

func (rem *Reminder) DurationUntilReminder() time.Duration {
	return rem.ReminderOccursAt().Sub(time.Now())
}

// struct to contain, and organize
// all pending reminders
type ReminderList struct {
	ReminderMap map[string]*Reminder
}

func (rl *ReminderList) AddReminder(rem *Reminder) {
	rl.ReminderMap[rem.ReminderIdentity()] = rem
}

func (rl *ReminderList) RemoveReminder(rem *Reminder) {
	delete(rl.ReminderMap, rem.ReminderIdentity())
}

func reminderStructFromCommand(recipient string, cmd string) (reminder *Reminder, err error) {

	r := &Reminder{Recipient: recipient, CreatedOn: time.Now()}
	r.ReminderIdentity()
	// see if cmd is empty
	cmd = strings.TrimSpace(cmd)
	if cmd == "" || recipient == "" {
		return nil, nil
	} else {

		args := strings.SplitAfterN(cmd, " ", 2)

		if len(args) == 1 {
			return nil, errors.New(INSUFFICENT_ARGS)
		} else {
			durationStr := strings.TrimSpace(args[0])
			reminderStr := strings.TrimSpace(args[1])

			// see if durationStr starts with any value except a digit
			firstChar := rune(durationStr[0])
			if (firstChar == '.') || unicode.IsDigit(firstChar) {
				r.Duration, err = time.ParseDuration(durationStr)
				if err != nil {
					return nil, err
				} else {
					r.Message = reminderStr
					return r, nil
				}
			} else {
				return nil, errors.New(fmt.Sprintf("Invalid duration value:[%s] supplied for argument.  Ignoring.", durationStr))
			}
		}
	}
}
