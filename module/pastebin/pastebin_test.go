package pastebin

import (
	// "fmt"
	"github.com/gamelost/bot3server/server"
	"log"
	"testing"
)

func TestPastebin2(t *testing.T) {

	svc := &PastebinService{PostURL: "http://greed.blackcore.com:7777", PostPath: "/documents"}
	response := &server.BotResponse{}

	response.Response = []string{"bas", "foo", "bar"}

	resp, err := svc.CreatePastebin(response.LinesAsByte())
	if err != nil {
		t.Fail()
	} else {
		log.Printf("Response is %s\n", resp)
	}
}
