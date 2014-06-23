package mongo

import (
	// "errors"
	// "encoding/json"
	"fmt"
	"github.com/gamelost/bot3server/server"
	// "io/ioutil"
	"log"
	// "math/rand"
	// "net/http"
	// "regexp"
	"strings"
	// "time"
	// "unicode"
	"labix.org/v2/mgo"
	// "labix.org/v2/mgo/bson"
	iniconf "code.google.com/p/goconf/conf"
	"launchpad.net/goyaml"
)

type MongoService struct {
	server.BotHandlerService
	ok      bool
	session *mgo.Session
	db      *mgo.Database
}

func (svc *MongoService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {

	var newSvc = &MongoService{}
	newSvc.BotHandlerService.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan

	newSvc.ok = false

	servers, err := newSvc.Config.GetString("mongo", "servers")

	if err != nil {
		log.Printf("Mongo: No server configured. Disabling")
		return newSvc
	}

	newSvc.session, err = mgo.Dial(servers)

	if err != nil {
		log.Printf("Unable to connect to Mongo: %v", err)
		return newSvc
	}

	db, err := newSvc.Config.GetString("mongo", "db")

	if err != nil {
		log.Printf("Mongo: No database configured, disabling.")
		return newSvc
	}

	newSvc.db = newSvc.session.DB(db)

	newSvc.ok = true

	return newSvc
}

func (svc *MongoService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	if !svc.ok {
		botResponse.SetSingleLineResponse("Mongo disabled")
		return
	}

	cmd := botRequest.Text()

	cmd = strings.TrimPrefix(cmd, "!mongo")
	cmd = strings.Trim(cmd, " ")
	resp := "Mongo is connected (I hope?)."

	if len(cmd) > 0 {
		splits := strings.SplitN(cmd, " ", 3)
		cmd := splits[0]
		var collection, arg string
		var marg map[string]interface{}
		var err error

		if len(splits) > 1 {
			collection = splits[1]
		}
		if len(splits) > 2 {
			arg = splits[2]
			err = goyaml.Unmarshal([]byte(arg), &marg)
			log.Printf("Unmarshaling '%s' to:", arg)
			for key, value := range marg {
				log.Println("Key:", key, "Value:", value)
			}
			if err != nil {
				log.Println("Unable to parse argument:", err)
				resp = fmt.Sprintf("Unable to parse argument: %+v", err)
			}
		}

		if err == nil {
			if cmd == "insert" {
				resp = svc.doInsert(collection, &marg)
			} else if cmd == "find" {
				resp = svc.doFind(collection, &marg)
			}
		}
	}
	botResponse.SetSingleLineResponse(resp)
	svc.PublishBotResponse(botResponse)
}

func (svc *MongoService) doInsert(collection string, arg *map[string]interface{}) string {
	c := svc.db.C(collection)
	err := c.Insert(*arg)
	if err == nil {
		return "Inserted?"
	} else {
		return fmt.Sprintf("Error on insert: %+v", err)
	}
}

func (svc *MongoService) doFind(collection string, arg *map[string]interface{}) string {
	c := svc.db.C(collection)
	results := c.Find(arg)

	if results == nil {
		return "Empty Resultset"
	} else {
		count, _ := results.Count()
		if count > 1 {
			return fmt.Sprintf("There were %d results for your find.", count)
		} else if count < 0 {
			return "No results."
		} else {
			var res interface{}
			err := results.One(&res)
			if err == nil {
				return fmt.Sprintf("%+v", res)
			} else {
				return fmt.Sprintf("Error while fetching result: %+v", err)
			}

		}
	}
}
