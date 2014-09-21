package server

import (
	iniconf "code.google.com/p/goconf/conf"
	"github.com/gamelost/bot3server/util"
	"time"
)

type BotHandler interface {
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

func (response *BotResponse) LinesAsByte() []byte {

	payload := make([]byte, 256)
	for _, value := range response.Response {
		payload = append(payload, value...)
		payload = append(payload, "\n"...)
	}
	return payload
}

func (request *BotRequest) RequestIsCommand() bool {

	return util.StringIsCommand(request.Text())
}

func (request *BotRequest) Command() string {

	return util.GetCommandFromString(request.Text())
}

func (request *BotRequest) LineTextWithoutCommand() string {

	if request.RequestIsCommand() {
		return util.TrimCommandFromString(request.Text(), request.Command())
	} else {
		return request.Text()
	}
}
