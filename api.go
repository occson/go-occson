package occson

import (
	"strings"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"crypto/rand"
	"io"
	"bytes"
)

func (workspace *Workspace) Download(url, passphrase string) []byte {
	url = strings.Replace(url, SCHEME, API, 1)

	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", "Token token=" + workspace.Token)
	req.Header.Set("Content-Type", "application/json")

	res, _ := client.Do(req)

	// TODO: Check response code, react accordingly
	// Assume 200 for now
	defer res.Body.Close()

	// 3. Get the body
	body, _ := ioutil.ReadAll(res.Body)

	// 4. Unmarshall into a struct
	response := Response{}
	_ = json.Unmarshal(body, &response)

	return ccsDecrypt(response.EncryptedContent, passphrase)
}

func (workspace *Workspace) Upload(url, passphrase, content string, force bool) bool {
	salt := make([]byte, 8)
    _, _ = io.ReadFull(rand.Reader, salt)

	ciph := ccsEncrypt(content, passphrase, string(salt))

	force_string := ""

	if force {
		force_string = "true"
	} else {
		force_string = "false"
	}

	request := Request{EncryptedContent: ciph, Force: force_string}

	jsonBody, _ := json.Marshal(request)
	client := http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	req.Header.Set("Authorization", "Token token="+ workspace.Token)
	req.Header.Set("Content-Type", "application/json")

	_, _ = client.Do(req)

	return true
}
