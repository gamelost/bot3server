package sleep

import (
	"fmt"
	"github.com/gamelost/bot3server/server"
	"time"
)

func Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	time.Sleep(10 * time.Second)
	botResponse.SetSingleLineResponse(fmt.Sprintf("Slept for ten seconds. Feeling refreshed %s!!", botRequest.Nick))
}
