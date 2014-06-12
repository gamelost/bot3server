package logger

import (
	//    "fmt"
	rdb "github.com/christopherhesse/rethinkgo"
	"github.com/gamelost/bot3server/server"
	"log"
)

var session *rdb.Session = nil

func init() {

	newSess, err := rdb.Connect("localhost:28015", "bot3")
	session = newSess
	if err != nil {
		log.Fatal("Failed to connect")
	}

	// try creating the log table in case it doesn't exist
	//err = rdb.DbCreate("bot3").Run(session).Exec()
	//err = rdb.TableCreate("log").Run(session).Exec()
	//if err != nil {
	//	log.Fatal("Failed to create table", err)
	//}
	//log.Println("Done setting up")
}

func Log(request *server.BotRequest) {

	err := rdb.Table("log").Insert(request).Run(session).Exec()
	if err != nil {
		log.Fatal("error occured", err)
	}
}
