package notehub

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type NotehubResponse struct {
	NoteID   string
	LongURL  string
	ShortURL string
	Status   *ResponseStatus
}

type ResponseStatus struct {
	Success bool
	Comment string
}

type NotehubService struct {
	PostURL            string
	NotehubCredentials *NotehubCredentials
}

type NotehubCredentials struct {
	PublisherId        string
	PublisherSecretKey string
}

func (cred *NotehubCredentials) isValidCredentials() bool {
	if cred.PublisherId == "" || cred.PublisherSecretKey == "" {
		return false
	}

	return true
}

func (nhs *NotehubService) CreateSignature(content []byte) ([]byte, error) {

	if nhs.NotehubCredentials.isValidCredentials() {
		hash := md5.New()
		io.WriteString(hash, nhs.NotehubCredentials.PublisherId)
		io.WriteString(hash, nhs.NotehubCredentials.PublisherSecretKey)
		io.WriteString(hash, string(content))
		return hash.Sum(nil), nil
	} else {
		return make([]byte, 0), errors.New("Missing required credentials for Notehub API service.")
	}
}

func (nhs *NotehubService) CreateDocument(content []byte) (*NotehubResponse, error) {

	// create MD5 hash
	hash, err := nhs.CreateSignature(content)
	if err != nil {
		return nil, err
	}

	values := make(url.Values)
	values.Set("note", string(content))
	values.Set("pid", nhs.NotehubCredentials.PublisherId)
	values.Set("signature", fmt.Sprintf("%x", hash))
	values.Set("version", "1.4")

	resp, err := http.PostForm(nhs.PostURL, values)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		notehubResponse := &NotehubResponse{}
		notehubResponse.Status = &ResponseStatus{}
		json.Unmarshal(body, notehubResponse)
		return notehubResponse, nil
	}
}
