package inconceivable

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"github.com/gamelost/bot3server/server"
)

type InconceivableService struct {
	server.BotHandlerService
}

func (svc *InconceivableService) NewService(config *iniconf.ConfigFile) server.BotHandler {
	newSvc := &InconceivableService{}
	newSvc.Config = config
	return newSvc
}

func (svc *InconceivableService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse(fmt.Sprintf("I do not think this word means what you think it means, %s!", botRequest.Nick))
	botResponse.ResponseType = "ACTION"
}
