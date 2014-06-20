package remindme

import (
	"errors"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"strings"
	"time"
	"unicode"
)

// set a max duration
const MAXDURATION = time.Hour * 24 * 7

// set a min duration
const MINDURATION = time.Second * 2

// string messages
const INSUFFICENT_ARGS = "Insufficent number of arguments provided.  Need to provide a duration and message."

type Reminder struct {
	Duration time.Duration
	NotifyAt time.Time
	Message  string
}

type RemindMeService struct {
	Reminders map[string]*Reminder
}

func (svc *RemindMeService) NewService() server.BotHandler {

	return &RemindMeService{}
}

func (svc *RemindMeService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	arg := botRequest.LineTextWithoutCommand()
	rem, err := HandleCommand(arg)

	if err != nil {
		botResponse.SetSingleLineResponse(fmt.Sprintf("Bloop. Your request could not be parsed: %s", err.Error()))
	} else {

		// nil reminder triggers status update instead
		if rem == nil {
			botResponse.SetSingleLineResponse(fmt.Sprintf("<placeholder for reminder summary>"))
			return
		} else if rem.Duration < 0 {
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s, only your mom would ask you to do something in the past. You're lame.", botRequest.Nick))
		} else if rem.Duration < MINDURATION {
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s, I dont work that fast!", botRequest.Nick))
		} else if rem.Duration > MAXDURATION {
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s, really? Maybe you should use a calendar instead.  Durations less than a week please.", botRequest.Nick))
		} else {
			time.Sleep(rem.Duration)
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s, you asked me to remind you: %s", botRequest.Nick, rem.Message))
		}
	}
}

func ReminderStructFromCommand(cmd string) (reminder *Reminder, err error) {

	r := &Reminder{}
	// see if cmd is empty
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
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

func HandleCommand(cmd string) (rem *Reminder, err error) {

	reminder, err := ReminderStructFromCommand(cmd)
	return reminder, err
}
