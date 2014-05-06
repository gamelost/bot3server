package dice

import "testing"
import "log"
import "fmt"

var svc *DiceService

func init() {
	svc = svc.NewService().(*DiceService)
}

func TestRoll(t *testing.T) {

	result := svc.RollDice("1d4")
	log.Printf(fmt.Sprintf("Result: %d, Rolls: %s", result.Result, result.Details))

	r2 := svc.RollDice("2d6;1d4")
	log.Printf(fmt.Sprintf("Result: %d, Rolls: %s", r2.Result, r2.Details))
}
