package s

import (
	iniconf "code.google.com/p/goconf/conf"
	"github.com/gamelost/bot3server/server"
	"strings"
	"fmt"
)

type SService struct {
	server.BotHandlerService
}

func (svc *SService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {
	var newSvc = &SService{}
	return newSvc
}

func (svc *SService) DispatchRequest(botRequest *server.BotRequest) {
	return
}

func (svc *SService) SubStringToStatement(sub string, written string) string {
	ary := strings.Split(sub[1:], "/")
	first = ary[0]
	second = ary[1]
	return
}
