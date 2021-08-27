package occson

import (
	"strings"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	aes "github.com/mervick/aes-everywhere/go/aes256"
	"errors"
)

const SCHEME = "ccs://"
const API = "https://api.occson.com/"

type Response struct {
	Id               string
	Path             string
	EncryptedContent string `json:"encrypted_content"`
	WorkspaceId      string `json:"workspace_id"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type Request struct {
	EncryptedContent string `json:"encrypted_content"`
	Force			 string   `json:"force"`
}

type Workspace struct {
	Token	string
}

func (workspace *Workspace) Download(url, passphrase string) ([]byte, error) {
	url = strings.Replace(url, SCHEME, API, 1)

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return []byte(""), err
	}

	req.Header.Set("Authorization", "Token token=" + workspace.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return []byte(""), err
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return []byte(""), errors.New(string(body))
	}

	response := Response{}
	err = json.Unmarshal(body, &response)

	if err != nil {
		return []byte(""), err
	}

	return []byte(aes.Decrypt(response.EncryptedContent, passphrase)), nil
}

func (workspace *Workspace) Upload(url, passphrase, content string, force bool) error {
	ciph := aes.Encrypt(content, passphrase)

	force_string := ""

	if force {
		force_string = "true"
	} else {
		force_string = "false"
	}

	request := Request{EncryptedContent: ciph, Force: force_string}

	jsonBody, err := json.Marshal(request)

	if err != nil {
		return err
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Token token="+ workspace.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(string(body))
	}

	return nil
}
