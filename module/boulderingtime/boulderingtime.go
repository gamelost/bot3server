package boulderingtime

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"github.com/gamelost/bot3server/server"
)

type BoulderingTimeService struct {
	server.BotHandlerService
}

func (svc *BoulderingTimeService) NewService(config *iniconf.ConfigFile) server.BotHandler {

	var newSvc = &BoulderingTimeService{}
	newSvc.Config = config
	return newSvc
}

func (svc *BoulderingTimeService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse(fmt.Sprintf("Its always bouldering time!"))
}
