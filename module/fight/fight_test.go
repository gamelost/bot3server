package fight

import "testing"
import "log"

var svc *FightService

func init() {
	svc = svc.NewService().(*FightService)
}

func TestRoll(t *testing.T) {

	result := svc.Fight("timmy", "paul")
	log.Printf("Result: %s", result)
}
