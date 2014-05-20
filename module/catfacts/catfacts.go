package catfacts

import (
	"github.com/gamelost/bot3server/server"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"regexp"
	"fmt"
)

type CatFactsResult struct {
	Success string   `json:"success"`
	Facts   []string `json:"facts"`
}

const baseUri = "http://catfacts-api.appspot.com/api/facts?number="

type CatFactsService struct{}

func (svc *CatFactsService) NewService() server.BotHandler {
	return &CatFactsService{}
}

func (svc *CatFactsService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {
	what := botRequest.RawLine.Text()
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
