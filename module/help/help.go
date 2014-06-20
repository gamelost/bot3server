package help

import (
	//	"fmt"
	iniconf "code.google.com/p/goconf/conf"
	"github.com/gamelost/bot3server/server"
)

type HelpService struct {
	server.BotHandlerService
}

func (svc *HelpService) NewService(config *iniconf.ConfigFile) server.BotHandler {

	var newSvc = &HelpService{}
	newSvc.Config = config
	return newSvc
}

func (svc *HelpService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse("i wont make you coffee and give a reach-around but you can ask the following: !remindme !fight !cah !inconceivable !slap !weather !forecast !nextwedding !zed")
}
