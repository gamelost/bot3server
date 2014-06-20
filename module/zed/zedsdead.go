package zed

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"github.com/gamelost/bot3server/server"
)

type ZedsDeadService struct {
	server.BotHandlerService
}

func (svc *ZedsDeadService) NewService(config *iniconf.ConfigFile) server.BotHandler {
	newSvc := &ZedsDeadService{}
	newSvc.Config = config
	return newSvc
}

func (svc *ZedsDeadService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse(fmt.Sprintf("Zed's dead baby. Zed's dead."))
}
