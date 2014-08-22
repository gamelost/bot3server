package help

import (
	//	"fmt"
	iniconf "code.google.com/p/goconf/conf"
	"github.com/gamelost/bot3server/server"
)

type HelpService struct {
	server.BotHandlerService
}

func (svc *HelpService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {

	var newSvc = &HelpService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	return newSvc
}

func (svc *HelpService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	botResponse.SetSingleLineResponse("i wont make you coffee and give a reach-around but you can ask the following: !remindme !fight !cah !inconceivable !slap !weather !forecast !nextwedding !nextbaby !zed !seen")
	svc.PublishBotResponse(botResponse)
}
