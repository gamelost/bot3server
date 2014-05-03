package inconceivable

import (
	"fmt"
	"github.com/gamelost/bot3server/server"
)

type InconceivableService struct{}

func (svc *InconceivableService) NewService() server.BotHandler {
	return &InconceivableService{}
}

func (svc *InconceivableService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse(fmt.Sprintf("I do not think this word means what you think it means, %s!", botRequest.RawLine.Nick))
	botResponse.ResponseType = "ACTION"
}
