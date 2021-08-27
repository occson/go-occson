# go-occson

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
