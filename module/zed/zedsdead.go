package zed

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"github.com/gamelost/bot3server/server"
)

type ZedsDeadService struct {
	server.BotHandlerService
}

func (svc *ZedsDeadService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {
	newSvc := &ZedsDeadService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	return newSvc
}

func (svc *ZedsDeadService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	botResponse.SetSingleLineResponse(fmt.Sprintf("Zed's dead baby. Zed's dead."))
	svc.PublishBotResponse(botResponse)
}
