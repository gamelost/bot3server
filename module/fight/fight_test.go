package fight

import "testing"
import "log"

func TestFight1(t *testing.T) {

	result := Fight("what do you think")
	
	if result != "Unable to parse fight command.  Please use 'vs' or 'vs.'" {
		t.Errorf("Should have passed back parse-fail message")
	}
}

func TestFight2(t *testing.T) {

	result := Fight("tim vs demislave")
	
	if result == "Unable to parse fight command.  Please use 'vs' or 'vs.'" {
		t.Errorf("Should have passed back parse-fail message")
	}

	log.Printf("%s", result)
}