package main

import (
	iniconf "code.google.com/p/goconf/conf"
	"encoding/json"
	"fmt"

	"github.com/alanjcfs/bot3server/module/cah"
	nsq "github.com/bitly/go-nsq"
	// "github.com/gamelost/bot3server/module/cah"
	"github.com/alanjcfs/bot3server/module/mongo"
	"github.com/gamelost/bot3server/module/catfacts"
	"github.com/gamelost/bot3server/module/dice"
	"github.com/gamelost/bot3server/module/fight"
	"github.com/gamelost/bot3server/module/help"
	"github.com/gamelost/bot3server/module/inconceivable"
	// "github.com/gamelost/bot3server/module/mongo"

	"github.com/gamelost/bot3server/module/nextwedding"
	"github.com/gamelost/bot3server/module/remindme"
	"github.com/gamelost/bot3server/module/slap"
	wuconditions "github.com/gamelost/bot3server/module/weather/conditions"
	wuforecast "github.com/gamelost/bot3server/module/weather/forecast"
	"github.com/gamelost/bot3server/module/zed"
	"github.com/gamelost/bot3server/server"
	"github.com/twinj/uuid"
	"labix.org/v2/mgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const DEFAULT_CONFIG_FILENAME = "bot3server.config"
const CONFIG_CAT_DEFAULT = "default"
const CONFIG_CAT_PLUGINS = "plugins"

// nsq specific constants
const CONFIG_CAT_NSQ = "nsq"
const CONFIG_BOT3SERVER_INPUT = "bot3server-input"
const CONFIG_BOT3SERVER_OUTPUT = "bot3server-output"
const CONFIG_OUTPUT_WRITER_ADDR = "output-writer-address"
const CONFIG_LOOKUPD_ADDR = "lookupd-address"
const TOPIC_MAIN = "main"

var conf *iniconf.ConfigFile

func main() {

	// the quit channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	config, err := iniconf.ReadConfigFile(DEFAULT_CONFIG_FILENAME)
	if err != nil {
		log.Fatal("Unable to read configuration file. Exiting now.")
	}

	bot3serverInput, _ := config.GetString(CONFIG_CAT_DEFAULT, CONFIG_BOT3SERVER_INPUT)
	outputWriterAddress, _ := config.GetString(CONFIG_CAT_NSQ, CONFIG_OUTPUT_WRITER_ADDR)
	lookupdAddress, _ := config.GetString(CONFIG_CAT_NSQ, CONFIG_LOOKUPD_ADDR)

	// set up listener instance
	incomingFromIRC, err := nsq.NewConsumer(bot3serverInput, TOPIC_MAIN, nsq.NewConfig())
	if err != nil {
		panic(err)
		sigChan <- syscall.SIGINT
	}

	// set up channels
	outgoingToNSQChan := make(chan *server.BotResponse)

	outputWriter, err := nsq.NewProducer(outputWriterAddress, nsq.NewConfig())
	if err != nil {
		panic(err)
		sigChan <- syscall.SIGINT
	}

	// set up heartbeat ticker
	heartbeatTicker := time.NewTicker(1 * time.Second)

	// initialize the handlers
	bot3server := &Bot3Server{Config: config, OutgoingChan: outgoingToNSQChan}
	bot3server.Initialize()

	incomingFromIRC.SetHandler(bot3server)
	incomingFromIRC.ConnectToNSQLookupd(lookupdAddress)

	go bot3server.HandleOutgoing(outgoingToNSQChan, heartbeatTicker.C, outputWriter, bot3server.UniqueID)

	log.Printf("Done starting up. UUID:[%s]. Waiting on quit signal.", bot3server.UniqueID.String)
	<-sigChan
}

type Bot3Server struct {
	Initialized  bool
	Config       *iniconf.ConfigFile
	Handlers     map[string]server.BotHandler
	OutgoingChan chan *server.BotResponse
	UniqueID     uuid.UUID
	MongoSession *mgo.Session
	MongoDB      *mgo.Database
}

func (bs *Bot3Server) Initialize() {
	// run only once
	if !bs.Initialized {

		bs.UniqueID = uuid.NewV1()
		bs.initServices()
		bs.Initialized = true
	}
}

func (bs *Bot3Server) AddHandler(key string, h server.BotHandler) {
	plugins, err := bs.Config.GetString(CONFIG_CAT_PLUGINS, "enabled")
	// If plugins string does not exist, assume that all plugins
	// are enabled.
	if err == nil {
		if !strings.Contains(" "+plugins+" ", " "+key+" ") {
			return
		}
	}
	bs.Handlers[key] = h
}

func (bs *Bot3Server) GetHandler(key string) server.BotHandler {
	return bs.Handlers[key]
}

func (bs *Bot3Server) initServices() error {

	bs.Handlers = make(map[string]server.BotHandler)

	// implement all services
	bs.AddHandler("fight", (new(fight.FightService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("cah", (new(cah.CahService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("mongo", (new(mongo.MongoService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("slap", (new(slap.SlapService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("inconceivable", (new(inconceivable.InconceivableService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("help", (new(help.HelpService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("remindme", (new(remindme.RemindMeService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("nextwedding", (new(nextwedding.NextWeddingService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("weather", (new(wuconditions.WeatherConditionsService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("forecast", (new(wuforecast.WeatherForecastService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("zed", (new(zed.ZedsDeadService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("dice", (new(dice.DiceService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("catfacts", (new(catfacts.CatFactsService)).NewService(bs.Config, bs.OutgoingChan))
	return nil
}

func (bs *Bot3Server) HandleMessage(message *nsq.Message) error {
	var req = &server.BotRequest{}
	json.Unmarshal(message.Body, req)

	err := bs.doInsert(req)
	if err != nil {
		return err
	}

	//	bs.IncomingChan <- message
	go bs.HandleIncoming(req)
	return nil
}

func (bs *Bot3Server) doInsert(req *server.BotRequest) error {
	if bs.MongoDB == nil {
		log.Printf("Setting up connection and inserting")
		bs.SetupMongoDBConnection()
	} else {
		log.Printf("Already connected to Mongo, inserting")
	}

	c := bs.MongoDB.C("chatlog")
	err := c.Insert(map[string]string{req.Nick: req.Text()})
	if err != nil {
		return err
	}
	return nil
}

func (bs *Bot3Server) SetupMongoDBConnection() error {
	// Connect to Mongo.
	servers, err := bs.Config.GetString("mongo", "servers")
	if err != nil {
		return err
	}

	bs.MongoSession, err = mgo.Dial(servers)
	if err != nil {
		return err
	}

	db, err := bs.Config.GetString("mongo", "db")
	if err != nil {
		return err
	} else {
		fmt.Println("Successfully obtained config from mongo")
	}

	bs.MongoDB = bs.MongoSession.DB(db)
	return nil
}

func (bs *Bot3Server) HandleOutgoing(outgoingToNSQChan chan *server.BotResponse, heartbeatTicker <-chan time.Time, outputWriter *nsq.Producer, serverID uuid.UUID) {

	for {
		select {
		case msg := <-outgoingToNSQChan:
			val, _ := json.Marshal(msg)
			outputWriter.Publish("bot3server-output", val)
			break
		case t := <-heartbeatTicker:
			hb := &server.Bot3ServerHeartbeat{ServerID: serverID.String(), Timestamp: t}
			val, _ := json.Marshal(hb)
			outputWriter.Publish("bot3server-heartbeat", val)
			break
		}
	}
}

func (bs *Bot3Server) HandleIncoming(botRequest *server.BotRequest) error {

	// since we dont want the entire botapi to crash if a module throws an error or panic
	// we'll trap it here and log/notify the owner.

	// mabye a good idea to automatically-disable the module if it panics too frequently

	// throwout malformed requests
	if len(botRequest.ChatText) < 1 {
		return nil
	}

	// deferred method to protect against panics
	defer func(rawline string) {
		// if module panics, log it and discard
		if x := recover(); x != nil {
			log.Printf("Trapped panic. Rawline is:[%s]: %v\n", rawline, x)
		}
	}(botRequest.Text())

	// check if command before processing
	if botRequest.RequestIsCommand() {

		command := botRequest.Command()

		handler := bs.GetHandler(command)

		if handler != nil {
			//log.Printf("Assigning handler for: %s\n", command)
			//botResponse := &server.BotResponse{Target: botRequest.Channel, Identifier: botRequest.Identifier}
			handler.DispatchRequest(botRequest)
		}
	}

	// return error object (nil for now)
	return nil
}
