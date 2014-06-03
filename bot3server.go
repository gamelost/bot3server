package main

import (
	"github.com/gamelost/bot3server/module/fight"
	// "github.com/gamelost/bot3server/module/logger"
	"github.com/gamelost/bot3server/server"
	//"github.com/gamelost/bot3server/module/helloworld"
	"github.com/gamelost/bot3server/module/help"
	"github.com/gamelost/bot3server/module/inconceivable"
	"github.com/gamelost/bot3server/module/slap"
	//"github.com/gamelost/bot3server/module/sleep"
	"github.com/gamelost/bot3server/module/cah"
	"github.com/gamelost/bot3server/module/nextwedding"
	// "github.com/gamelost/bot3server/module/panic"
	"github.com/gamelost/bot3server/module/dice"
	"github.com/gamelost/bot3server/module/remindme"
	"github.com/gamelost/bot3server/module/zed"
	//"github.com/gamelost/bot3server/module/testcommands"
	"github.com/gamelost/bot3server/module/boulderingtime"
	"github.com/gamelost/bot3server/module/catfacts"
	wuconditions "github.com/gamelost/bot3server/module/weather/conditions"
	wuforecast "github.com/gamelost/bot3server/module/weather/forecast"
	//"github.com/gamelost/bot3server/module/howlongbeforeicanquityouintel"
	iniconf "code.google.com/p/goconf/conf"
	"encoding/json"
	nsq "github.com/gamelost/go-nsq"
	"github.com/twinj/uuid"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DEFAULT_CONFIG_FILENAME = "bot3server.config"
const CONFIG_CAT_DEFAULT = "default"

// nsq specific constants
const CONFIG_CAT_NSQ = "nsq"
const CONFIG_BOT3SERVER_INPUT = "bot3server-input"
const CONFIG_BOT3SERVER_OUTPUT = "bot3server-output"
const CONFIG_OUTPUT_WRITER_ADDR = "output-writer-address"
const CONFIG_LOOKUPD_ADDR = "lookupd-address"
const TOPIC_MAIN = "main"

func main() {

	// create unique UUID on startup
	newUUID := uuid.NewV1()

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
	incomingFromIRC, err := nsq.NewReader(bot3serverInput, TOPIC_MAIN)
	if err != nil {
		panic(err)
		sigChan <- syscall.SIGINT
	}

	// set up channels
	outgoingToNSQChan := make(chan *server.BotResponse)

	outputWriter := nsq.NewWriter(outputWriterAddress)

	// set up heartbeat ticker
	heartbeatTicker := time.NewTicker(1 * time.Second)

	// initialize the handlers
	botApp := &BotApp{Config: config, OutgoingChan: outgoingToNSQChan, UniqueID: newUUID}
	botApp.initServices()

	incomingFromIRC.AddHandler(botApp)
	incomingFromIRC.ConnectToLookupd(lookupdAddress)

	go HandleOutgoingToNSQ(outgoingToNSQChan, heartbeatTicker.C, outputWriter, newUUID)

	log.Printf("Done starting up. UUID:[%s]. Waiting on quit signal.", newUUID.String())
	<-sigChan
}

func HandleOutgoingToNSQ(outgoingToNSQChan chan *server.BotResponse, heartbeatTicker <-chan time.Time, outputWriter *nsq.Writer, serverID uuid.UUID) {

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

type BotApp struct {
	Config       *iniconf.ConfigFile
	Handlers     map[string]server.BotHandler
	OutgoingChan chan *server.BotResponse
	UniqueID     uuid.UUID
}

func (ba *BotApp) AddHandler(key string, h server.BotHandler) {
	ba.Handlers[key] = h
}

func (ba *BotApp) GetHandler(key string) server.BotHandler {
	return ba.Handlers[key]
}

func (ba *BotApp) initServices() error {

	ba.Handlers = make(map[string]server.BotHandler)

	// implement all services
	ba.AddHandler("fight", (new(fight.FightService)).NewService())
	ba.AddHandler("cah", (new(cah.CahService)).NewService())
	ba.AddHandler("slap", (new(slap.SlapService)).NewService())
	ba.AddHandler("inconceivable", (new(inconceivable.InconceivableService)).NewService())
	ba.AddHandler("help", (new(help.HelpService)).NewService())
	ba.AddHandler("remindme", (new(remindme.RemindMeService)).NewService())
	ba.AddHandler("nextwedding", (new(nextwedding.NextWeddingService)).NewService())
	ba.AddHandler("weather", (new(wuconditions.WeatherConditionsService)).NewService())
	ba.AddHandler("forecast", (new(wuforecast.WeatherForecastService)).NewService())
	ba.AddHandler("zed", (new(zed.ZedsDeadService)).NewService())
	ba.AddHandler("boulderingtime", (new(boulderingtime.BoulderingTimeService)).NewService())
	ba.AddHandler("dice", (new(dice.DiceService)).NewService())
	ba.AddHandler("catfacts", (new(catfacts.CatFactsService)).NewService())
	return nil
}

func (ba *BotApp) HandleMessage(message *nsq.Message) error {

	//	ba.IncomingChan <- message
	var req = &server.BotRequest{}
	json.Unmarshal(message.Body, req)
	go ba.HandleIncoming(req)
	return nil
}

func (ba *BotApp) HandleIncoming(botRequest *server.BotRequest) error {

	// since we dont want the entire botapi to crash if a module throws an error or panic
	// we'll trap it here and log/notify the owner.

	// mabye a good idea to automatically-disable the module if it panics too frequently

	// throwout malformed requests
	if botRequest.RawLine == nil {
		return nil
	}

	// deferred method to protect against panics
	defer func(rawline string) {
		// if module panics, log it and discard
		if x := recover(); x != nil {
			log.Printf("Trapped panic. Rawline is:[%s]: %v\n", rawline, x)
		}
	}(botRequest.RawLine.Text())

	// log all lines
	// logger.Log(botRequest)

	// check if command before processing
	if botRequest.RequestIsCommand() {

		command := botRequest.Command()

		handler := ba.GetHandler(command)

		if handler != nil {
			//log.Printf("Assigning handler for: %s\n", command)
			botResponse := &server.BotResponse{Target: botRequest.RawLine.Target()}
			handler.Handle(botRequest, botResponse)
			ba.OutgoingChan <- botResponse
		}
	}

	// return error object (nil for now)
	return nil
}
