package remindme

import (
	"fmt"
	"github.com/gamelost/bot3server/server"
	"strings"
	"time"
)

type RemindMeService struct{}

func (svc *RemindMeService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	arg := botRequest.LineTextWithoutCommand()
	resp, dur, err := HandleCommand(arg)

	if err != nil {
		botResponse.SetSingleLineResponse("Could not understand your request.  Should be in [duration] [message] format.")
	} else {

		if dur < 0 {
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s, only your mom would ask you to do something in the past.  You're lame.", botRequest.RawLine.Nick))
		} else {
			time.Sleep(dur)
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s, you asked me to remind you: %s", botRequest.RawLine.Nick, resp))
		}
	}
}

func HandleCommand(cmd string) (response string, dur time.Duration, err error) {

	cmd = strings.TrimSpace(cmd)
	args := strings.SplitAfterN(cmd, " ", 2)

	dur, err = time.ParseDuration(strings.TrimSpace(args[0]))
	if err != nil {
		return
	} else {
		response = strings.TrimSpace(args[1])
		return
	}
}
