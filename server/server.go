package server

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"strings"
)

type BotHandler interface {
	NewService() BotHandler
	Handle(botRequest *BotRequest, botResponse *BotResponse)
}

const (
	PRIVMSG = "PRIVMSG"
	NOTICE  = "NOTICE"
	ACTION  = "ACTION"
)

type BotRequest struct {
	// raw data struct (from the bot3)
	RawLine *irc.Line
}

type BotResponse struct {
	ResponseType string
	Target       string
	Response     []string
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

func (request *BotRequest) RequestIsCommand() bool {

	return stringIsCommand(request.RawLine.Text())
}

func (request *BotRequest) Command() string {

	return getCommandFromString(request.RawLine.Text())
}

func (request *BotRequest) LineTextWithoutCommand() string {

	if request.RequestIsCommand() {
		lineTxt := strings.TrimPrefix(request.RawLine.Text(), fmt.Sprintf("!%s", request.Command()))
		return strings.TrimSpace(lineTxt)
	} else {
		return request.RawLine.Text()
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
