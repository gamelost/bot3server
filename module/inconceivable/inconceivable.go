package inconceivable

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"github.com/gamelost/bot3server/server"
)

type InconceivableService struct {
	server.BotHandlerService
}

func (svc *InconceivableService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {
	newSvc := &InconceivableService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	return newSvc
}

func (svc *InconceivableService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	botResponse.SetSingleLineResponse(fmt.Sprintf("I do not think this word means what you think it means, %s!", botRequest.Nick))
	botResponse.ResponseType = "ACTION"
	svc.PublishBotResponse(botResponse)
}
