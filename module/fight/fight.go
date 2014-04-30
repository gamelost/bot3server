package fight

import (
	"errors"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"math/rand"
	"strings"
	"time"
)

var rng *rand.Rand
var fightMethodSlice []string
var fightArenaSlice []string

type FightService struct {
}

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	fightMethodSlice = []string{"tickles to death", "pummels", "quarters", "garrotes", "butchers", "obliterates", "tears apart limb by limb", "annihilates", "rampages past", "dismembers", "kneecaps", "uses force lightning to crispy-critter", "gets blown out of the sky by", "executes a well-timed Harai goshi on", "smothers"}
	fightArenaSlice = []string{"in a gentlemanly game of chess", "in a fight to the pain", "on the dark side of the moon", "in the mens restroom", "in the ladies restroom", "in a barroom brawl", "in a slapfest", "with dull flaming scimitars", "on the planet Hoth", "with elephant foreskins filled with brie"}
}

func (svc *FightService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	fighterOne, fighterTwo, err := parseInput(botRequest.RawLine.Text())

	if err != nil {
		botResponse.SetSingleLineResponse("Unable to parse fight command.  Please use 'vs' or 'vs.'")
	} else {
		randVal := rng.Intn(2)

		if randVal == 0 {
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s %s %s %s.", fighterOne, randomFightVerb(), fighterTwo, randomFightArena()))
		} else {
			botResponse.SetSingleLineResponse(fmt.Sprintf("%s %s %s %s.", fighterTwo, randomFightVerb(), fighterOne, randomFightArena()))
		}
	}
}

func randomFightVerb() string {

	randVal := rng.Intn(len(fightMethodSlice))
	method := fightMethodSlice[randVal]
	return method
}

func randomFightArena() string {

	randVal := rng.Intn(len(fightArenaSlice))
	method := fightArenaSlice[randVal]
	return method
}

func parseInput(input string) (fighterOne string, fighterTwo string, err error) {

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
