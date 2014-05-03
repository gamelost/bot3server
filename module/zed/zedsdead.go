package zed

import (
	"fmt"
	"github.com/gamelost/bot3server/server"
)

type ZedsDeadService struct{}

func (svc *ZedsDeadService) NewService() server.BotHandler {
	return &ZedsDeadService{}
}

func (svc *ZedsDeadService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse(fmt.Sprintf("Zed's dead baby.  Zed's dead."))
}
