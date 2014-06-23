package dice

import (
	iniconf "code.google.com/p/goconf/conf"
	"encoding/json"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"io/ioutil"
	"log"
	"net/http"
	url "net/url"
	"strings"
)

type DiceService struct {
	server.BotHandlerService
}

type DiceResult struct {
	Error       string `json:"error"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	Description string `json:"description"`
}

func (svc *DiceService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {
	newSvc := &DiceService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	return newSvc
}

func (svc *DiceService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	log.Println("Received Handle for !dice")
	result := svc.RollDice(parseInput(botRequest.Text()))
	strResp := fmt.Sprintf("%s: your roll:[%s] result:[%s], description:[%s]", botRequest.Nick, result.Input, result.Output, result.Description)
	botResponse.SetSingleLineResponse(strResp)
	svc.PublishBotResponse(botResponse)
}

func (svc *DiceService) RollDice(diceString string) *DiceResult {

	log.Println("Rolling dice")
	// download the CAH json file
	var escapedInput = url.QueryEscape(diceString)
	var urlStr = "http://lethalcode.net:8080/roll?src=" + escapedInput
	log.Println(urlStr)
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Output: %s\n", body)
	result := &DiceResult{}

	err = json.Unmarshal(body, result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON: %s\n", result)
	return result
}

func parseInput(input string) string {

	input = strings.TrimPrefix(input, "!dice")
	diceRequest := strings.TrimSpace(input)
	return diceRequest
}
