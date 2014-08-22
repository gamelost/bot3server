package seen

import (
	iniconf "code.google.com/p/goconf/conf"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

type SeenService struct {
	server.BotHandlerService
	MongoSession *mgo.Session
	MongoDB      *mgo.Database
}

func (svc *SeenService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {

	var newSvc = &SeenService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	newSvc.setupMongoDBConnection()
	return newSvc
}

func (svc *SeenService) DispatchRequest(botRequest *server.BotRequest) {

	botResponse := svc.CreateBotResponse(botRequest)
	strInput := parseInput(botRequest.Text())
	seenRes, err := DurationSinceLastMessage(strInput, botRequest.Channel, svc.MongoDB)

	if err != nil {
		botResponse.SetSingleLineResponse(fmt.Sprintf("Could not find user: <%s> in logs.  No idea when last seen.", strInput))
	} else {
		str := fmt.Sprintf("Last seen: <%s> %s, saying: '%s'", seenRes.Nick, humanizeDuration(seenRes.DurationSinceMessage()), seenRes.Text)
		botResponse.SetSingleLineResponse(str)
	}

	svc.PublishBotResponse(botResponse)
}

func parseInput(input string) string {

	input = strings.TrimPrefix(input, "!seen")
	input = strings.Trim(input, " ")

	return input
}

func humanizeDuration(dur time.Duration) string {

	if dur.Seconds() < 1 {
		return "less than a second ago"
	} else if dur.Seconds() < 60 {
		return fmt.Sprintf("%1.f second(s) ago", dur.Seconds())
	} else if dur.Minutes() < 60 {
		if dur.Minutes() == 1 {
			return "about a minute ago."
		} else {
			return fmt.Sprintf("%1.f minutes ago", dur.Minutes())
		}
	} else if dur.Hours() < 24 {
		if dur.Hours() == 1 {
			return "about a hour ago."
		} else {
			return fmt.Sprintf("%1.f hours ago", dur.Hours())
		}
	} else {
		return fmt.Sprintf("%1.f days ago", dur.Hours()/24)
	}
}

// Connect to Mongo reading from the configuration file
func (svc *SeenService) setupMongoDBConnection() error {

	servers, err := svc.Config.GetString("mongo", "servers")
	if err != nil {
		panic(err)
	}

	svc.MongoSession, err = mgo.Dial(servers)
	if err != nil {
		panic(err)
	}

	db, err := svc.Config.GetString("mongo", "db")
	if err != nil {
		panic(err)
	} else {
		log.Println("Successfully obtained config from mongo")
	}

	svc.MongoDB = svc.MongoSession.DB(db)
	return nil
}

type SeenResult struct {
	Nick    string        `bson:"nick"`
	Channel string        `bson:"channel"`
	Id      bson.ObjectId `bson:"_id"`
	Text    string        `bson:"chattext"`
}

func (res *SeenResult) DurationSinceMessage() time.Duration {
	dur := time.Now().Sub(res.Id.Time())
	return dur
}

func DurationSinceLastMessage(username string, channel string, database *mgo.Database) (sr *SeenResult, err error) {

	c := database.C("chatlog")
	result := &SeenResult{}
	err = c.Find(bson.M{"nick": bson.RegEx{Pattern: "^" + username, Options: "i"}, "channel": channel}).Sort("-_id").One(result)

	if err != nil {
		// fmt.Println(err.Error())
		return nil, err
	}

	return result, nil
}
