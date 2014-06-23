package server

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"strings"
	"time"
)

type BotHandler interface {
	NewService(config *iniconf.ConfigFile, outgoingToIRCChan chan *BotResponse) BotHandler
	DispatchRequest(botRequest *BotRequest)
}

func (bhs *BotHandlerService) CreateBotResponse(botRequest *BotRequest) *BotResponse {
	return &BotResponse{Target: botRequest.Channel, Identifier: botRequest.Identifier}
}

func (bhs *BotHandlerService) PublishBotResponse(botResponse *BotResponse) {
	bhs.PublishToIRCChan <- botResponse
}

type BotHandlerService struct {
	PublishToIRCChan chan *BotResponse
	Config           *iniconf.ConfigFile
}

const (
	PRIVMSG = "PRIVMSG"
	NOTICE  = "NOTICE"
	ACTION  = "ACTION"
)

type BotRequest struct {
	// raw data struct (from the bot3)
	Identifier string
	Nick       string
	Channel    string
	ChatText   string
}

func (req *BotRequest) Text() string {
	return req.ChatText
}

type BotResponse struct {
	Identifier   string
	ResponseType string
	Target       string
	Response     []string
}

type Bot3ServerHeartbeat struct {
	ServerID  string
	Timestamp time.Time
}

func (response *BotResponse) IsMultiLineResponse() bool {

	if len(response.Response) > 1 {
		return true
	} else {
		return false
	}
}

func (response *BotResponse) SingleLineResponse() string {

	return response.Response[0]
}

func (response *BotResponse) SetSingleLineResponse(rstr string) {

	responseArr := []string{rstr}
	response.Response = responseArr
}

func (response *BotResponse) SetMultipleLineResponse(rstr []string) {

	response.Response = rstr
}

func (request *BotRequest) RequestIsCommand() bool {

	return stringIsCommand(request.Text())
}

func (request *BotRequest) Command() string {

	return getCommandFromString(request.Text())
}

func (request *BotRequest) LineTextWithoutCommand() string {

	if request.RequestIsCommand() {
		lineTxt := strings.TrimPrefix(request.Text(), fmt.Sprintf("!%s", request.Command()))
		return strings.TrimSpace(lineTxt)
	} else {
		return request.Text()
	}
}

func stringIsCommand(rawstring string) bool {

	if strings.HasPrefix(rawstring, "!") {
		return true
	} else {
		return false
	}
}

func getCommandFromString(rawstring string) string {

	var commandStr = ""

	if stringIsCommand(rawstring) {
		index := strings.Index(rawstring, " ")
		if index > 0 {
			commandStr = rawstring[1:index]
		} else if index == -1 {
			// no space, so return entirety of string except first char
			commandStr = rawstring[1:len(rawstring)]
		}
	}

	return commandStr
}
