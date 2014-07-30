package main

import (
	"fmt"
	"testing"

	"github.com/gamelost/bot3server/server"
)

var botServer *Bot3Server
var botReq = server.BotRequest{}

func Init() {
	botServer = new(Bot3Server)
	botServer.SetupMongoDBConnection()
}

func TestSanity(t *testing.T) {
	const a, b = 2, 2
	if a+b == 4 {
		fmt.Println("All working")
	}
}

func TestAddHandler(t *testing.T) {
}

func TestDoInsert(t *testing.T) {
	// botServer.DoInsert(botReq)
	// botServer.MongoDB.C("chatlog").Find(bson.M{"nick": "dormiens"})
}
