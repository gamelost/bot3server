package seen

import (
	// "fmt"
	"gopkg.in/mgo.v2"
	"log"
	"testing"
)

func TestStructFromCommand1(t *testing.T) {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	database := session.DB("gamelost")
	seenRes, err := DurationSinceLastMessage("dormiens", "#tpsreports", database)
	log.Printf("Last seen: <%s> about %s ago, saying:'%s'", seenRes.Nick, humanizeDuration(seenRes.DurationSinceMessage()), seenRes.Nick)
}
