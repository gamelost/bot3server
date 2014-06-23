package slap

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"strings"
)

type SlapService struct {
	server.BotHandlerService
}

func (svc *SlapService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {
	newSvc := &SlapService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	return newSvc
}

func (svc *SlapService) DispatchRequest(botRequest *server.BotRequest) {

	victim := parseInput(botRequest.Text())

	br := svc.CreateBotResponse(botRequest)
	br.SetSingleLineResponse(fmt.Sprintf("slaps %s with a fine Corinthian leather glove hand-sewn from the remains of a thousand Buick Cordobas.", victim))
	br.ResponseType = server.ACTION
	svc.PublishBotResponse(br)
}

func parseInput(input string) string {

	input = strings.TrimPrefix(input, "!slap")
	victim := strings.TrimSpace(input)
	return victim
}
