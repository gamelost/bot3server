package helloworld

import (
	"github.com/gamelost/bot3server/server"
)

func Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	resp := []string{"Hello, " + botRequest.RawLine.Nick + "!", "How are you today?"}
	botResponse.Response = resp
}
