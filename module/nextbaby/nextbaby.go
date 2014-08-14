package nextbaby

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"time"
)

type NextBabyService struct {
	server.BotHandlerService
}

func (svc *NextBabyService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {
	newSvc := &NextBabyService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	return newSvc
}

func (svc *NextBabyService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	botResponse.SetSingleLineResponse(durationToDate())
	svc.PublishBotResponse(botResponse)
}

func durationToDate() string {

	nowDate := time.Now()
	loc, _ := time.LoadLocation("America/Los_Angeles")
	weddingDate := time.Date(2015, time.February, 24, 12, 0, 0, 0, loc)

	if nowDate.After(weddingDate) {
		return "Ah mon. It be too late, has da baby popped out yet?!"
	} else {

		duration := weddingDate.Sub(nowDate)
		durationInMinutes := int(duration.Minutes())
		durationDays := durationInMinutes / (24 * 60)
		durationHours := (durationInMinutes - (durationDays * 60 * 24)) / 60
		durationMinutes := durationInMinutes - ((durationDays * 60 * 24) + (durationHours * 60))

		return fmt.Sprintf("sonOfAshburn rises in %d days, %d hours and %d minutes! Have you shopped at Baby-R-Us yet?!", durationDays, durationHours, durationMinutes)
	}
}
