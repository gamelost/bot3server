package notehub

import (
	// "fmt"
	// "github.com/gamelost/bot3server/server"
	"log"
	"testing"
)

func TestNotehub(t *testing.T) {

	svc := &NotehubService{PostURL: "http://www.notehub.org/api/note"}

	creds := &NotehubCredentials{PublisherId: "bot3server", PublisherSecretKey: "2af52476393a83e2529fdbd881ddb34a"}
	svc.NotehubCredentials = creds
	resp, err := svc.CreateDocument([]byte("_Freedom!_\n\nWar!!\n\nSalvation!\nEmpire of the sun!"))
	if err != nil {
		log.Printf("Error occured:", err)
	}
	log.Printf("Response: %v", resp.LongURL)
}
