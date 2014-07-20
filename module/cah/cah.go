package cah

import (
	// "errors"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	iniconf "code.google.com/p/goconf/conf"
	"github.com/gamelost/bot3server/server"
)

// source URL for CAH card templates
const CAH_SOURCE_URL = "https://raw.githubusercontent.com/gamelost/bot3server/master/module/cah/cah-cards-standard.json"

type CahService struct {
	server.BotHandlerService
	RandomNG          *rand.Rand
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

func (svc *CahService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {

	var newSvc = &CahService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	newSvc.CahCardCollection = &CahCardCollection{}

	// set up the rng
	newSvc.RandomNG = rand.New(rand.NewSource(time.Now().UnixNano()))

	// download the CAH json file
	resp, err := http.Get(CAH_SOURCE_URL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, newSvc.CahCardCollection)

	return newSvc
}

func (svc *CahService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	strInput := parseInput(botRequest.Text())
	if len(strInput) > 0 {
		botResponse.SetSingleLineResponse(svc.RandomCahMessageWithArgument(strInput))
	} else {
		botResponse.SetSingleLineResponse(svc.RandomCahMessage())
	}

	svc.PublishBotResponse(botResponse)
}

func (svc *CahService) RandomCard() *CahCard {

	randVal := svc.RandomNG.Intn(svc.CahCardCollection.CardCount())
	return svc.CahCardCollection.GetCardAt(randVal)
}

func (svc *CahService) RandomQuestionCard() *CahCard {

	cCard := svc.RandomCard()
	for {
		if cCard.CardType == "Q" {
			break
		} else {
			cCard = svc.RandomCard()
		}
	}
	return cCard
}

func (svc *CahService) RandomOneAnswerQuestionCard() *CahCard {

	cCard := svc.RandomQuestionCard()
	for {
		if cCard.NumAnswers == 1 {
			break
		} else {
			cCard = svc.RandomQuestionCard()
		}
	}
	return cCard
}

func (svc *CahService) RandomAnswerCard() *CahCard {

	cCard := svc.RandomCard()
	for {
		if cCard.CardType == "A" {
			break
		} else {
			cCard = svc.RandomCard()
		}
	}
	return cCard
}

func (svc *CahService) MessageFromQuestionAndAnswers(questionStr string, answers []string) string {

	var finalStr string
	substrings := strings.Split(questionStr, "_")
	if len(substrings) < 2 {
		finalStr = questionStr + " " + convertToStandaloneAnswer(answers[0])
	} else {

		ansCounter := 0
		for _, value := range substrings {
			finalStr += value
			if ansCounter < len(answers) {
				finalStr += convertToInlineAnswer(answers[ansCounter])
				ansCounter++
			}
		}
	}

	return finalStr
}

func (svc *CahService) RandomCahMessage() string {

	qCard := svc.RandomQuestionCard()

	// find out how many answers we need
	numAnswers := qCard.NumAnswers

	// queue up all needed answer cards
	var answers = make([]string, numAnswers)
	for i := 0; i < int(numAnswers); i++ {
		answers[i] = svc.RandomAnswerCard().Text
	}

	return svc.MessageFromQuestionAndAnswers(qCard.Text, answers)
}

func (svc *CahService) RandomCahAnswerMessage(s string) string {
	card := svc.RandomAnswerCard()
	return convertToInlineAnswer(card.Text)
}

func (svc *CahService) RandomCahMessageWithArgument(argStr string) string {

	log.Printf("We got %s", argStr)
	if strings.Contains(argStr, "__") {
		re, _ := regexp.Compile("__+")
		ret := re.ReplaceAllStringFunc(argStr, svc.RandomCahAnswerMessage)
		return ret
	} else {
		qCard := svc.RandomOneAnswerQuestionCard()

		// queue up all needed answer cards
		var answers = make([]string, 1)
		answers[0] = argStr

		return svc.MessageFromQuestionAndAnswers(qCard.Text, answers)
	}
}

func (ccc CahCardCollection) CardCount() int {
	return len(ccc)
}

func (ccc CahCardCollection) GetCardAt(cardLoc int) *CahCard {
	return &CahCard{CardType: ccc[cardLoc].CardType, Id: ccc[cardLoc].Id, NumAnswers: ccc[cardLoc].NumAnswers, Text: ccc[cardLoc].Text}
}

func parseInput(input string) string {

	input = strings.TrimPrefix(input, "!cah")
	input = strings.Trim(input, " ")

	return input
}

func convertToInlineAnswer(orig string) string {

	runes := []rune(orig)
	runes[0] = unicode.ToLower(runes[0])
	orig = string(runes)
	return strings.TrimRight(orig, ".")
}

func convertToStandaloneAnswer(orig string) string {

	runes := []rune(orig)
	runes[0] = unicode.ToUpper(runes[0])

	// if no punctuation, add a period
	lastChar := runes[len(runes)-1]
	if lastChar == '.' || lastChar == '!' || lastChar == '?' {
		return string(runes)
	} else {
		return string(runes) + "."
	}
}
