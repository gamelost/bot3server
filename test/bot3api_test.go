package bot3api

import (
	irc "github.com/fluffle/goirc/client"
	"github.com/gamelost/bot3server/server"
	"log"
	"net/rpc"
	"testing"
)

const TEST_RPC_URL = "http://localhost:8888/rpc"
const RPC_ENCODING = "application/json"

func BenchmarkRPC(b *testing.B) {

	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// close client before function exit
	defer client.Close()

	line := &irc.Line{}

	botRequest := &server.BotRequest{RawLine: line}
	botResponse := &server.BotResponse{}

	err = client.Call("BotService.Handle", botRequest, &botResponse)
	if err != nil {
		log.Fatal("#gamelost", "Unable to comply, botserver seems to have gone away... %s", err)
	}

	log.Println("Done with call.")
}
