package occson

import (
	"context"
	"crypto/tls"
	"fmt"
	aes "github.com/mervick/aes-everywhere/go/aes256"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, s.Close
}

var encryptedContent = aes.Encrypt("test", "test")

var okResponse = fmt.Sprintf(`{
	"id": "1",
	"path": "/test",
	"encrypted_content": "%s",
	"workspace_id": "1",
	"created_at": "1",
	"updated_at": "1"
}`, encryptedContent)

var errorResponse = `{"message":"fail", "status":404}`

func TestDocumentDownloadSuccess(t *testing.T) {
	uri := "ccs://test.toml"
	token := "deadbeef"
	passphrase := "test"

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Token token=deadbeef", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		assert.Equal(t, "api.occson.com", r.Host)
		assert.Equal(t, "/test.toml", r.URL.Path)

		w.Write([]byte(okResponse))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	document := Document{Uri: uri, Token: token, Passphrase: passphrase, httpClient: httpClient}

	result, err := document.Download()

	assert.Nil(t, err)
	assert.Equal(t, []byte("test"), result)
}

func TestDocumentDownloadFailure(t *testing.T) {
	uri := "ccs://test.toml"
	token := "deadbeef"
	passphrase := "test"

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Token token=deadbeef", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		assert.Equal(t, "api.occson.com", r.Host)
		assert.Equal(t, "/test.toml", r.URL.Path)

		http.Error(w, errorResponse, 404)
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	document := Document{Uri: uri, Token: token, Passphrase: passphrase, httpClient: httpClient}

	result, err := document.Download()

	assert.Equal(t, err.Error(), errorResponse+"\n")
	assert.Equal(t, []byte(""), result)
}

func TestDocumentUploadSuccess(t *testing.T) {
	uri := "ccs://test.toml"
	token := "deadbeef"
	passphrase := "test"

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Token token=deadbeef", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		assert.Equal(t, "api.occson.com", r.Host)
		assert.Equal(t, "/test.toml", r.URL.Path)

		body, _ := ioutil.ReadAll(r.Body)

		assert.Contains(t, string(body), `"encrypted_content":`)
		assert.Contains(t, string(body), `"force":"false"`)

		w.Write([]byte(okResponse))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	document := Document{Uri: uri, Token: token, Passphrase: passphrase, httpClient: httpClient}

	err := document.Upload("test", false)

	assert.Nil(t, err)
}

func TestDocumentForcedUploadSuccess(t *testing.T) {
	uri := "ccs://test.toml"
	token := "deadbeef"
	passphrase := "test"

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Token token=deadbeef", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		assert.Equal(t, "api.occson.com", r.Host)
		assert.Equal(t, "/test.toml", r.URL.Path)

		body, _ := ioutil.ReadAll(r.Body)

		assert.Contains(t, string(body), `"encrypted_content":`)
		assert.Contains(t, string(body), `"force":"true"`)

		w.Write([]byte(okResponse))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	document := Document{Uri: uri, Token: token, Passphrase: passphrase, httpClient: httpClient}

	err := document.Upload("test", true)

	assert.Nil(t, err)
}

func TestDocumentUploadFailure(t *testing.T) {
	uri := "ccs://test.toml"
	token := "deadbeef"
	passphrase := "test"

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Token token=deadbeef", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		http.Error(w, errorResponse, 403)
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	document := Document{Uri: uri, Token: token, Passphrase: passphrase, httpClient: httpClient}

	err := document.Upload("test", false)

	assert.Equal(t, err.Error(), errorResponse+"\n")
}

func TestUrlConversion(t *testing.T) {
	uri := "ccs://test.toml"
	doc := NewDocument(uri, "", "")

	assert.Equal(t, doc.url(), "https://api.occson.com/test.toml")
}
