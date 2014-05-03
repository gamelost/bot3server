package slap

import (
	"fmt"
	"github.com/gamelost/bot3server/server"
	"strings"
)

type SlapService struct {
}

func (svc *SlapService) NewService() server.BotHandler {
	return &SlapService{}
}

func (svc *SlapService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	victim := parseInput(botRequest.RawLine.Text())
	botResponse.SetSingleLineResponse(fmt.Sprintf("slaps %s with a fine Corinthian leather glove hand-sewn from the remains of a thousand Buick Cordobas.", victim))
	botResponse.ResponseType = server.ACTION
}

func parseInput(input string) string {

	input = strings.TrimPrefix(input, "!slap")
	victim := strings.TrimSpace(input)
	return victim
}
