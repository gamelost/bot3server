package help

import (
	//	"fmt"
	"github.com/gamelost/bot3server/server"
)

type HelpService struct{}

func (svc *HelpService) NewService() server.BotHandler {
	return &HelpService{}
}

func (svc *HelpService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse("i wont make you coffee and give a reach-around but you can ask the following: !remindme !fight !cah !inconceivable !slap !sleep !weather !forecast !nextwedding !zed")
}
