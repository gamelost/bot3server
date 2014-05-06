package dice

import (
	"encoding/json"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type DiceService struct {
}

type DiceResult struct {
	Result       int64  `json:"result"`
	Details      string `json:"details"`
	Code         string `json:"code"`
	Timestamp    int64  `json:"timestamp"`
	Illustration string `json:"illustration"`
}

func (svc *DiceService) NewService() server.BotHandler {
	return &DiceService{}
}

func (svc *DiceService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	log.Println("Received Handle for !dice")
	result := svc.RollDice(parseInput(botRequest.RawLine.Text()))
	strResp := fmt.Sprintf("%s: your roll (%s) result: %d, Rolls: %s", botRequest.RawLine.Nick, result.Code, result.Result, result.Details)
	botResponse.SetSingleLineResponse(strResp)
}

func (svc *DiceService) RollDice(diceString string) *DiceResult {

	log.Println("Rolling dice")
	// download the CAH json file
	resp, err := http.Get("http://rolz.org/api/?" + diceString + ".json")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
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
