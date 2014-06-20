package catfacts

import (
	iniconf "code.google.com/p/goconf/conf"
	"encoding/json"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

type CatFactsResult struct {
	Success string   `json:"success"`
	Facts   []string `json:"facts"`
}

const baseUri = "http://catfacts-api.appspot.com/api/facts?number="

type CatFactsService struct {
	server.BotHandlerService
}

func (svc *CatFactsService) NewService(config *iniconf.ConfigFile) server.BotHandler {
	newSvc := &CatFactsService{}
	newSvc.Config = config
	return newSvc
}

func (svc *CatFactsService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {
	what := botRequest.Text()
	number := getBoundedIntegerFromInput(what, 1, 10)
	result := svc.CatFactsApi(number)
	botResponse.SetMultipleLineResponse(result.Facts)
}

func (svc *CatFactsService) CatFactsApi(number int) *CatFactsResult {

	resp, err := http.Get(fmt.Sprintf("%s%d", baseUri, number))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	result := &CatFactsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		panic(err)
	}

	return result
}

func getBoundedIntegerFromInput(raw string, min int, max int) int {
	r, _ := regexp.Compile("!\\w+\\s+(\\d+)\\s*$")
	match := r.FindStringSubmatch(raw)
	if len(match) != 2 {
		return min
	}

	number, err := strconv.Atoi(match[1])
	if err != nil { // shouldn't happen!
		return min
	}
	if number > max {
		return max
	}
	if number < min {
		return min
	}
	return number
}
