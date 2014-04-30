package remindme

import "testing"
import "log"

func TestHandleCommand(t *testing.T) {

	str, _, err := HandleCommand("1m coffee is ready")
	
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	
	log.Println(str)
}