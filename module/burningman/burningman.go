package burningman

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"time"
)

type BurningManService struct {
	server.BotHandlerService
}

func (svc *BurningManService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {
	newSvc := &BurningManService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	return newSvc
}

func (svc *BurningManService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	botResponse.SetSingleLineResponse(durationToBurningMan())
	svc.PublishBotResponse(botResponse)
}

func durationToBurningMan() string {

	nowDate := time.Now()
	loc, _ := time.LoadLocation("America/Los_Angeles")
	partyBeginsDate := time.Date(2015, time.August, 31, 10, 0, 0, 0, loc)
	manBurnsDate := time.Date(2015, time.September, 5, 22, 0, 0, 0, loc)

	if nowDate.After(manBurnsDate) {
		return "Ashes to ashes.  The man's done burnt!"
	} else if nowDate.Before(manBurnsDate) && nowDate.After(partyBeginsDate) {
		return "Party aint over yet!  The man still standing!"
	} else {
		partyBeginsDuration := partyBeginsDate.Sub(nowDate)
		partyBeginsDurationInMinutes := int(partyBeginsDuration.Minutes())
		partyBeginsDurationDays := partyBeginsDurationInMinutes / (24 * 60)
		partyBeginsDurationHours := (partyBeginsDurationInMinutes - (partyBeginsDurationDays * 60 * 24)) / 60
		partyBeginsDurationMinutes := partyBeginsDurationInMinutes - ((partyBeginsDurationDays * 60 * 24) + (partyBeginsDurationHours * 60))
		//
		manBurnsDuration := manBurnsDate.Sub(partyBeginsDate)
		manBurnsDurationHours := int(manBurnsDuration.Hours())

		return fmt.Sprintf("Only %d days, %d hours and %d minutes left till the big party on the playa! And the man burns %d hours after!", partyBeginsDurationDays, partyBeginsDurationHours, partyBeginsDurationMinutes, manBurnsDurationHours)
	}
}
