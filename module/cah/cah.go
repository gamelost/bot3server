package cah

import (
	// "errors"
	"encoding/json"
	//"fmt"
	"github.com/gamelost/bot3server/server"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unicode"
)

var rng *rand.Rand
var cahCardCollection *CahCardCollection

type CahService struct {
	CahCardCollection *CahCardCollection
}

type CahCard struct {
	Id         int64
	CardType   string
	NumAnswers int64
	Text       string
}

type CahCardCollection []struct {
	CardType   string `json:"cardType"`
	Expansion  string `json:"expansion"`
	Id         int64  `json:"id"`
	NumAnswers int64  `json:"numAnswers"`
	Text       string `json:"text"`
}

func (ccc CahCardCollection) CardCount() int {
	return len(ccc)
}

func (ccc CahCardCollection) RandomCard() *CahCard {
	randVal := rng.Intn(len(ccc))
	return &CahCard{CardType: ccc[randVal].CardType, Id: ccc[randVal].Id, NumAnswers: ccc[randVal].NumAnswers, Text: ccc[randVal].Text}
}

func (ccc CahCardCollection) RandomQuestionCard() *CahCard {

	cCard := ccc.RandomCard()
	for {
		if cCard.CardType == "Q" {
			break
		} else {
			cCard = ccc.RandomCard()
		}
	}
	return cCard
}

func (ccc CahCardCollection) RandomOneAnswerQuestionCard() *CahCard {

	cCard := ccc.RandomCard()
	for {
		if cCard.CardType == "Q" && cCard.NumAnswers == 1 {
			break
		} else {
			cCard = ccc.RandomCard()
		}
	}
	return cCard
}

func (ccc CahCardCollection) RandomAnswerCard() *CahCard {

	cCard := ccc.RandomCard()
	for {
		if cCard.CardType == "A" {
			break
		} else {
			cCard = ccc.RandomCard()
		}
	}
	return cCard
}

func (ccc CahCardCollection) RandomCahMessage() string {

	var finalStr string
	qCard := ccc.RandomQuestionCard()

	// find out how many answers we need
	numAnswers := qCard.NumAnswers

	// queue up all needed answer cards
	var ansCards = make([]*CahCard, numAnswers)
	for i := 0; i < int(numAnswers); i++ {
		ansCards[i] = ccc.RandomAnswerCard()
	}

	substrings := strings.Split(qCard.Text, "_")
	if len(substrings) < 2 {
		finalStr = qCard.Text + " " + ansCards[0].Text
	} else {

		ansCounter := 0
		for _, value := range substrings {
			finalStr += value
			if ansCounter < len(ansCards) {
				finalStr += sanitizeAnswerText(ansCards[ansCounter].Text)
				ansCounter++
			}
		}
	}

	return finalStr
}

func (ccc CahCardCollection) RandomCahMessageWithArgument(arg string) string {

	var finalStr string
	qCard := ccc.RandomOneAnswerQuestionCard()

	// queue up all needed answer cards
	var ansCards = make([]*CahCard, 1)
	ansCards[0] = &CahCard{Text: arg}

	substrings := strings.Split(qCard.Text, "_")
	if len(substrings) < 2 {
		finalStr = qCard.Text + " " + ansCards[0].Text
	} else {

		ansCounter := 0
		for _, value := range substrings {
			finalStr += value
			if ansCounter < len(ansCards) {
				finalStr += sanitizeAnswerText(ansCards[ansCounter].Text)
				ansCounter++
			}
		}
	}

	return finalStr
}

func parseInput(input string) string {

	input = strings.TrimPrefix(input, "!cah")
	input = strings.Trim(input, " ")

	return input
}

func sanitizeAnswerText(orig string) string {

	runes := []rune(orig)
	runes[0] = unicode.ToLower(runes[0])
	orig = string(runes)
	return strings.TrimRight(orig, ".")
}

func init() {

	// set up the rng
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	// download the CAH json file
	// https://raw.githubusercontent.com/nodanaonlyzuul/against-humanity/master/source/cards.json
	resp, err := http.Get("https://raw.githubusercontent.com/nodanaonlyzuul/against-humanity/master/source/cards.json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	cahCardCollection = &CahCardCollection{}
	json.Unmarshal(body, cahCardCollection)

	//log.Printf("Message: %s", cahCardCollection.RandomCahMessageWithArgument("Timmy"))
	//log.Printf("Message2: %s", cahCardCollection.RandomCahMessage())

	log.Println("Done with init()")
}

func (svc *CahService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	strInput := parseInput(botRequest.RawLine.Text())
	if len(strInput) > 0 {
		botResponse.SetSingleLineResponse(cahCardCollection.RandomCahMessageWithArgument(strInput))
	} else {
		botResponse.SetSingleLineResponse(cahCardCollection.RandomCahMessage())
	}

}
