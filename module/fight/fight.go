package fight

import (
	iniconf "code.google.com/p/goconf/conf"
	"errors"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"math/rand"
	"strings"
	"time"
)

type FightService struct {
	server.BotHandlerService
	RandomNG     *rand.Rand
	FightMethods []string
	FightArenas  []string
}

func (svc *FightService) NewService(config *iniconf.ConfigFile) server.BotHandler {

	var newSvc = &FightService{}
	newSvc.Config = config
	newSvc.RandomNG = rand.New(rand.NewSource(time.Now().UnixNano()))
	newSvc.FightMethods = []string{"tickles to death", "pummels", "quarters", "garrotes", "butchers", "obliterates", "tears apart limb by limb", "annihilates", "rampages past", "dismembers", "kneecaps", "uses force lightning to crispy-critter", "gets blown out of the sky by", "executes a well-timed Harai goshi on", "smothers"}
	newSvc.FightArenas = []string{"in a gentlemanly game of chess", "in a fight to the pain", "on the dark side of the moon", "in the mens restroom", "in the ladies restroom", "in a barroom brawl", "in a slapfest", "with dull flaming scimitars", "on the planet Hoth", "with elephant foreskins filled with brie"}
	return newSvc
}

func (svc *FightService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	fighterOne, fighterTwo, err := svc.ParseInput(botRequest.Text())

	if err != nil {
		botResponse.SetSingleLineResponse("Unable to parse fight command.  Please use 'vs' or 'vs.'")
	} else {
		botResponse.SetSingleLineResponse(svc.Fight(fighterOne, fighterTwo))
	}
}

func (svc *FightService) Fight(fighterOne string, fighterTwo string) string {

	randVal := svc.RandomNG.Intn(2)

	if randVal == 0 {
		return fmt.Sprintf("%s %s %s %s.", fighterOne, svc.RandomFightVerb(), fighterTwo, svc.RandomFightArena())
	} else {
		return fmt.Sprintf("%s %s %s %s.", fighterTwo, svc.RandomFightVerb(), fighterOne, svc.RandomFightArena())
	}
}

func (svc *FightService) RandomFightVerb() string {

	randVal := svc.RandomNG.Intn(len(svc.FightMethods))
	method := svc.FightMethods[randVal]
	return method
}

func (svc *FightService) RandomFightArena() string {

	randVal := svc.RandomNG.Intn(len(svc.FightArenas))
	method := svc.FightArenas[randVal]
	return method
}

func (svc *FightService) ParseInput(input string) (fighterOne string, fighterTwo string, err error) {

	input = strings.TrimPrefix(input, "!fight ")

	fighters := strings.SplitAfterN(input, "vs", 2)

	if len(fighters) == 2 {

		fighterOne = strings.TrimSuffix(fighters[0], "vs")
		fighterTwo = strings.TrimPrefix(fighters[1], ".")

		fighterOne = strings.TrimSpace(fighterOne)
		fighterTwo = strings.TrimSpace(fighterTwo)
	} else {
		err := errors.New("number of fighters is not two")
		return "", "", err
	}

	return
}
