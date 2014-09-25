package main

import (
	iniconf "code.google.com/p/goconf/conf"
	"encoding/json"
	"fmt"
	nsq "github.com/bitly/go-nsq"
	"github.com/gamelost/bot3server/module/cah"
	"github.com/gamelost/bot3server/module/catfacts"
	"github.com/gamelost/bot3server/module/dice"
	"github.com/gamelost/bot3server/module/fight"
	"github.com/gamelost/bot3server/module/help"
	"github.com/gamelost/bot3server/module/inconceivable"
	"github.com/gamelost/bot3server/module/mongo"
	"github.com/gamelost/bot3server/module/nextbaby"
	"github.com/gamelost/bot3server/module/nextwedding"
	"github.com/gamelost/bot3server/module/pastebin"
	"github.com/gamelost/bot3server/module/remindme"
	"github.com/gamelost/bot3server/module/seen"
	"github.com/gamelost/bot3server/module/slap"
	"github.com/gamelost/bot3server/module/stats"
	wuconditions "github.com/gamelost/bot3server/module/weather/conditions"
	wuforecast "github.com/gamelost/bot3server/module/weather/forecast"
	"github.com/gamelost/bot3server/module/wq"
	"github.com/gamelost/bot3server/module/zed"
	"github.com/gamelost/bot3server/server"
	"github.com/twinj/uuid"
	"gopkg.in/mgo.v2"
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

// pastebin constants
const PASTEBIN_HOST = "pastebin-host"
const PASTEBIN_PATH = "pastebin-path"
const PASTEBIN_MAXLINES = "pastebin-maxlines"

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

	incomingFromIRC.AddHandler(bot3server)
	incomingFromIRC.ConnectToNSQLookupd(lookupdAddress)

	go bot3server.HandleOutgoing(outgoingToNSQChan, heartbeatTicker.C, outputWriter, bot3server.UniqueID)

	log.Printf("Done starting up. UUID:[%s]. Waiting on quit signal.", bot3server.UniqueID.String)
	<-sigChan
}

type Bot3Server struct {
	Initialized      bool
	Config           *iniconf.ConfigFile
	Handlers         map[string]server.BotHandler
	OutgoingChan     chan *server.BotResponse
	UniqueID         uuid.UUID
	MongoSession     *mgo.Session
	MongoDB          *mgo.Database
	PastebinService  *pastebin.PastebinService
	PastebinMaxLines int
}

func (bs *Bot3Server) Initialize() {
	// run only once
	if !bs.Initialized {

		bs.UniqueID = uuid.NewV1()
		bs.initServices()
		bs.SetupMongoDBConnection()
		bs.SetupPastebinService()
		bs.Initialized = true
	}
}

func (bs *Bot3Server) SetupPastebinService() {

	pbHost, _ := bs.Config.GetString("pastebin", PASTEBIN_HOST)
	pbPath, _ := bs.Config.GetString("pastebin", PASTEBIN_PATH)
	bs.PastebinMaxLines, _ = bs.Config.GetInt("pastebin", PASTEBIN_MAXLINES)

	bs.PastebinService = &pastebin.PastebinService{PostURL: pbHost, PostPath: pbPath}
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
	bs.AddHandler("remindme", (remindme.NewRemindMeService(bs.Config, bs.OutgoingChan)))
	bs.AddHandler("nextwedding", (new(nextwedding.NextWeddingService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("nextbaby", (new(nextbaby.NextBabyService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("weather", (new(wuconditions.WeatherConditionsService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("forecast", (new(wuforecast.WeatherForecastService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("wq", (new(wikiquote.WikiQuoteService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("zed", (new(zed.ZedsDeadService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("dice", (new(dice.DiceService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("catfacts", (new(catfacts.CatFactsService)).NewService(bs.Config, bs.OutgoingChan))
	bs.AddHandler("stats", (stats.NewStatsService(bs.Config, bs.OutgoingChan)))
	bs.AddHandler("seen", (new(seen.SeenService)).NewService(bs.Config, bs.OutgoingChan))
	return nil
}

func (bs *Bot3Server) HandleMessage(message *nsq.Message) error {

	var req = &server.BotRequest{}
	json.Unmarshal(message.Body, req)

	err := bs.DoInsert(req)
	if err != nil {
		return err
	}

	//	bs.IncomingChan <- message
	go bs.HandleIncoming(req)
	return nil
}

func (bs *Bot3Server) DoInsert(req *server.BotRequest) error {
	c := bs.MongoDB.C("chatlog")
	err := c.Insert(req)
	if err != nil {
		return err
	}
	return nil
}

// Connect to Mongo reading from the configuration file
func (bs *Bot3Server) SetupMongoDBConnection() error {

	servers, err := bs.Config.GetString("mongo", "servers")
	if err != nil {
		panic(err)
	}

	bs.MongoSession, err = mgo.Dial(servers)
	if err != nil {
		panic(err)
	}

	db, err := bs.Config.GetString("mongo", "db")
	if err != nil {
		panic(err)
	} else {
		log.Println("Successfully obtained config from mongo")
	}

	bs.MongoDB = bs.MongoSession.DB(db)
	return nil
}

func (bs *Bot3Server) HandleOutgoing(outgoingToNSQChan chan *server.BotResponse, heartbeatTicker <-chan time.Time, outputWriter *nsq.Producer, serverID uuid.UUID) {

	for {
		select {
		case msg := <-outgoingToNSQChan:
			bs.PreProcessAndPublishOutgoing(msg, outputWriter)
			break
		case t := <-heartbeatTicker:
			hb := &server.Bot3ServerHeartbeat{ServerID: serverID.String(), Timestamp: t}
			val, _ := json.Marshal(hb)
			outputWriter.Publish("bot3server-heartbeat", val)
			break
		}
	}
}

func (bs *Bot3Server) PreProcessAndPublishOutgoing(msg *server.BotResponse, outputWriter *nsq.Producer) {

	if len(msg.Response) > bs.PastebinMaxLines {
		totalLineLength := len(msg.Response)
		resp, _ := bs.PastebinService.CreatePastebin(msg.LinesAsByte())
		msg.Response = msg.Response[0:bs.PastebinMaxLines]

		resp = fmt.Sprintf("( ...remaining %d lines clipped. view all content at: %s )", (totalLineLength - bs.PastebinMaxLines), resp)
		msg.Response = append(msg.Response, resp)
	}

	val, _ := json.Marshal(msg)
	outputWriter.Publish("bot3server-output", val)
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
	} else {
		log.Println("Unable to handle incoming request: %v", botRequest)
	}
	// return error object (nil for now)
	return nil
}
