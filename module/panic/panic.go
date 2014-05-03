package panic

import (
	"github.com/gamelost/bot3server/server"
)

/* a simple service that throws panic when called.  used for
testing the ability for the api to receive the panic, recover and not end */
type PanicService struct{}

func (svc *PanicService) NewService() server.BotHandler {
	return &PanicService{}
}

func (svc *PanicService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {
	panic("The world is ending! Panic!")
}
