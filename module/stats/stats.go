package stats

import (
	iniconf "code.google.com/p/goconf/conf"
	// "encoding/json"
	// "errors"
	// "fmt"
	"github.com/gamelost/bot3server/server"
	// "github.com/twinj/uuid"
	// "io/ioutil"
	// "os"
	// "strings"
	// "time"
	// "unicode"
	//"log"
)

func NewStatsService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) *StatsService {
	newSvc := &StatsService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	return newSvc
}

type StatsService struct {
	server.BotHandlerService
}

func (svc *StatsService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	botResponse.SetSingleLineResponse("stats module")

	svc.PublishBotResponse(botResponse)
}
