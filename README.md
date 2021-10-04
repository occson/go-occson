# go-occson

![master CI status](https://github.com/occson/go-occson/actions/workflows/test.yml/badge.svg?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/occson/go-occson.svg)](https://pkg.go.dev/github.com/occson/go-occson)

This package provides a client for the API of [occson.com](https://occson.com) - a configuration control system.

## Downloading a document's contents

```go
import (
	occson "github.com/occson/go-occson"
)

func main() {
	// Not sure where to get these? Check out occson.com!
	uri := "ccs://golang-test.toml"
	token := "decafc0ffeebad"
	passphrase := "deadbeef"

	// Sets up the document struct using a helper
	doc := occson.NewDocument(uri, token, passphrase)

	// Performs the actual request and decryption
	decrypted, err := doc.Download()

	if err != nil {
		panic(err)
	}

	// Prints out plaintext of the document
	fmt.Println(string(decrypted))
}
```

## Uploading new encrypted contents

```go
import (
	occson "github.com/occson/go-occson"
)

func main() {
	uri := "ccs://golang-test.toml"
	token := "decafc0ffeebad"
	passphrase := "deadbeef"


	// Sets up the document struct using a helper
	doc := occson.NewDocument(uri, token, passphrase)

	// Our new plaintext contents
	blob := `
		[config]
		param = "some_param"
	`

	// Performs the encryption and upload of ciphertext
	doc.Upload(blob, true)
}

```
