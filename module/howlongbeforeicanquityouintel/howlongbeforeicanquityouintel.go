package howlongbeforeicanquityouintel

import (
	"fmt"
	"github.com/gamelost/bot3server/server"
	"time"
)

type HowLongBeforeICanQuitYouIntelService struct{}

func (svc *HowLongBeforeICanQuitYouIntelService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	botResponse.SetSingleLineResponse(durationTo())
}

func durationTo() string {

	nowDate := time.Now()
	loc, _ := time.LoadLocation("America/Los_Angeles")
	weddingDate := time.Date(2013, time.July, 8, 16, 0, 0, 0, loc)
	duration := weddingDate.Sub(nowDate)

	durationInMinutes := int(duration.Minutes())

	durationDays := durationInMinutes / (24 * 60)
	durationHours := (durationInMinutes - (durationDays * 60 * 24)) / 60
	durationMinutes := durationInMinutes - ((durationDays * 60 * 24) + (durationHours * 60))

	durationStr := fmt.Sprintf("Only %d days, %d hours and %d minutes left for MacCoaster at Intel.  He's gonna miss TFS terribly!", durationDays, durationHours, durationMinutes)

	return durationStr
}
