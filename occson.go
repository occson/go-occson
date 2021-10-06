package occson

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	aes "github.com/mervick/aes-everywhere/go/aes256"
)

// Defines the scheme which denotes a CCS URL.
const SCHEME = "occson://"

// API endpoint for Occson.
const API = "https://api.occson.com/"

type Document struct {
	// URI of the document. Can and should begin with the "occson://" schema.
	Uri string
	// Auth token of the appropriate workspace.
	Token string
	// Document passphrase.
	Passphrase string

	httpClient *http.Client
}

type Response struct {
	// Requested document's internal ID.
	Id string
	// Requested document's path.
	Path string
	// Requested document's encrypted content (before decryption).
	EncryptedContent string `json:"encrypted_content"`
	// Document's workspace internal ID.
	WorkspaceId string `json:"workspace_id"`
	// Document's creation time, in ISO8601 format.
	CreatedAt string `json:"created_at"`
	// Document's last update time, in ISO8601 format.
	UpdatedAt string `json:"updated_at"`
}

type Request struct {
	// Encrypted content for upload.
	EncryptedContent string `json:"encrypted_content"`
	// Whether the document should be overwritten, even if it exists.
	Force string `json:"force"`
}

// Helper function to create a document struct quickly.
func NewDocument(uri, token, passphrase string) Document {
	doc := Document{Uri: uri, Token: token, Passphrase: passphrase, httpClient: &http.Client{}}

	return doc
}

// Downloads the given document from Occson, returning its decrypted contents. Authentication
// uses the Token field, decryption uses the Passphrase field.
func (doc *Document) Download() ([]byte, error) {
	req, err := http.NewRequest("GET", doc.url(), nil)

	if err != nil {
		return []byte(""), err
	}

	req.Header.Set("Authorization", "Token token="+doc.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := doc.httpClient.Do(req)

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

	return []byte(aes.Decrypt(response.EncryptedContent, doc.Passphrase)), nil
}

// Encrypts and uploads the given content, optionally overwriting the contents
// already in Occson. Authentication uses the Token field, encryption uses the Passphrase field.
func (doc *Document) Upload(content string, force bool) error {
	ciph := aes.Encrypt(content, doc.Passphrase)

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

	req, err := http.NewRequest("POST", doc.url(), bytes.NewBuffer(jsonBody))

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Token token="+doc.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := doc.httpClient.Do(req)

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

func (doc *Document) url() string {
	return strings.Replace(doc.Uri, SCHEME, API, 1)
}
