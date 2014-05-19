package dice

import "testing"
import "log"
import "fmt"

var svc *DiceService

func init() {
	svc = svc.NewService().(*DiceService)
}

func TestRoll(t *testing.T) {

	result := svc.RollDice("(+ 1 3)")
	log.Printf(fmt.Sprintf("Result: %s, Rolls: %s", result.Input, result.Output))

	r2 := svc.RollDice("(+ 1 4)")
	log.Printf(fmt.Sprintf("Result: %s, Rolls: %s", r2.Input, r2.Output))
}
