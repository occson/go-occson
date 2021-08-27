# go-occson

[![Build Status](https://app.travis-ci.com/paweljw/go-occson.svg?branch=master)](https://app.travis-ci.com/paweljw/go-occson)
[![Go Reference](https://pkg.go.dev/badge/github.com/paweljw/go-occson.svg)](https://pkg.go.dev/github.com/paweljw/go-occson)

```go
package main

import (
	occson "github.com/paweljw/go-occson"
)

func main() {
	url := "ccs://golang-test.toml"
	token := "decafc0ffeebad"
	passphrase := "deadbeef"
	
	doc := occson.NewDocument(uri, token, passphrase)

	decrypted := doc.Download()
	fmt.Println(string(decrypted))

	blob := `[config]

param = "some_param"`

	doc.Upload(blob, true)
}
```
