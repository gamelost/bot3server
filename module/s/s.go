package s

import (
	iniconf "code.google.com/p/goconf/conf"
	"github.com/gamelost/bot3server/server"
)

type SService struct {
	server.BotHandlerService
}

func (svc *SService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {
	var newSvc = &SService{}
	return newSvc
}

func (svc *SService) DispatchRequest(botRequest *server.BotRequest) {
	return
}

func (svc *SService) SubStringToStatement(sub string, written string) string {
	return "I had a package golden - I lost it in the sand"
}
