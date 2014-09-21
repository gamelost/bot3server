package pastebin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PastebinKeyResponse struct {
	Key     string
	Message string
}

type PastebinService struct {
	PostURL  string
	PostPath string
}

func (pbs *PastebinService) CreatePastebin(content []byte) (string, error) {

	log.Printf("Creatin content for %s", content)
	buf := bytes.NewBuffer(content)
	resp, err := http.Post(pbs.PostURL+pbs.PostPath, "text/plain", buf)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error occured %s", err.Error())
		}
		log.Printf("Response is %s", body)
		keyStruct := &PastebinKeyResponse{}
		json.Unmarshal(body, keyStruct)

		url := fmt.Sprintf("%s/%s", pbs.PostURL, keyStruct.Key)
		return url, nil
	}
}
