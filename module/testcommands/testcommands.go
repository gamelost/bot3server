package testcommands

import (
	"fmt"
	"github.com/gamelost/bot3server/server"
)

func Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	resp := make([]string, 10)
	for i := 0; i < len(resp); i++ {
		resp[i] = fmt.Sprintf("Line %d", i)
	}

	botResponse.Response = resp
}
