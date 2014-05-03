package boulderingtime

import (
	"fmt"
	"github.com/gamelost/bot3server/server"
)

type BoulderingTimeService struct{}

func (svc *BoulderingTimeService) NewService() server.BotHandler {

	var newSvc = &BoulderingTimeService{}
	return newSvc
}

func (svc *BoulderingTimeService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse(fmt.Sprintf("Its always bouldering time!"))
}
