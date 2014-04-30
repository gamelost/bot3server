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
	"github.com/gamelost/bot3server/module/nextwedding"
	//"github.com/gamelost/bot3server/module/panic"
	"github.com/gamelost/bot3server/module/cah"
	"github.com/gamelost/bot3server/module/remindme"
	"github.com/gamelost/bot3server/module/zed"
	//"github.com/gamelost/bot3server/module/testcommands"
	"github.com/gamelost/bot3server/module/boulderingtime"
	wuconditions "github.com/gamelost/bot3server/module/weather/conditions"
	wuforecast "github.com/gamelost/bot3server/module/weather/forecast"
	//"github.com/gamelost/bot3server/module/howlongbeforeicanquityouintel"
	iniconf "code.google.com/p/goconf/conf"
	"encoding/json"
	"fmt"
	nsq "github.com/bitly/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// the quit channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	config, err := iniconf.ReadConfigFile("bot3api.config")
	if err != nil {
		log.Fatal("Unable to read configuration file. Exiting now.")
	}

	bot3apiInput, _ := config.GetString("default", "bot3api-input")

	// set up listener instance
	incomingFromIRC, err := nsq.NewReader(bot3apiInput, "main")
	if err != nil {
		panic(err)
		sigChan <- syscall.SIGINT
	}

	// set up channels
	incomingFromNSQChan := make(chan *nsq.Message)
	outgoingToNSQChan := make(chan *server.BotResponse)

	writer := nsq.NewWriter("127.0.0.1:4150")

	// initialize the handlers
	botApp := &BotApp{Config: config, IncomingChan: incomingFromNSQChan, OutgoingChan: outgoingToNSQChan, NSQOutputWriter: writer}
	botApp.initServices()

	incomingFromIRC.AddHandler(botApp)
	incomingFromIRC.ConnectToLookupd("127.0.0.1:4161")

	go botApp.DoBotService(incomingFromNSQChan, outgoingToNSQChan)

	<-sigChan
}

type BotApp struct {
	Config          *iniconf.ConfigFile
	Handlers        map[string]server.BotHandler
	IncomingChan    chan *nsq.Message
	OutgoingChan    chan *server.BotResponse
	NSQOutputWriter *nsq.Writer
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
	ba.AddHandler("fight", new(fight.FightService))
	ba.AddHandler("cah", new(cah.CahService))
	ba.AddHandler("slap", new(slap.SlapService))
	ba.AddHandler("inconceivable", new(inconceivable.InconceivableService))
	ba.AddHandler("help", new(help.HelpService))
	ba.AddHandler("remindme", new(remindme.RemindMeService))
	ba.AddHandler("nextwedding", new(nextwedding.NextWeddingService))
	ba.AddHandler("weather", new(wuconditions.WeatherConditionsService))
	ba.AddHandler("forecast", new(wuforecast.WeatherForecastService))
	//ba.AddHandler("panic", new(panic.PanicService))
	ba.AddHandler("zed", new(zed.ZedsDeadService))
	ba.AddHandler("boulderingtime", new(boulderingtime.BoulderingTimeService))
	return nil
}

func (ba *BotApp) DoBotService(incomingFromNSQChan chan *nsq.Message, outgoingToNSQChan chan *server.BotResponse) {

	for {
		select {
		case msg := <-incomingFromNSQChan:
			fmt.Printf("Received incoming message from NSQ bus: %s\n", msg)
			break
		case msg := <-outgoingToNSQChan:
			//fmt.Printf("Received outgoing message for NSQ bus: %s\n", msg)
			val, _ := json.Marshal(msg)
			ba.NSQOutputWriter.Publish("bot3api-output", val)
			break
		}
	}
}

func (ba *BotApp) HandleMessage(message *nsq.Message) error {

	ba.IncomingChan <- message
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
	fmt.Printf("Line is %s", botRequest)
	if botRequest.RawLine == nil {
		return nil
	}

	// deferred method to protect against panics
	defer func(rawline string) {
		// if module panics, log it and discard
		if x := recover(); x != nil {
			log.Printf("Trapped panic. Rawline is:[%s]: %v", rawline, x)
		}
	}(botRequest.RawLine.Text())

	// log all lines
	// logger.Log(botRequest)

	// check if command before processing
	if botRequest.RequestIsCommand() {

		command := botRequest.Command()

		handler := ba.GetHandler(command)

		if handler != nil {
			fmt.Printf("Handler for: %s\n", command)

			botResponse := &server.BotResponse{Target: botRequest.RawLine.Target()}
			handler.Handle(botRequest, botResponse)
			ba.OutgoingChan <- botResponse
		}
	}

	// return error object (nil for now)
	return nil
}
