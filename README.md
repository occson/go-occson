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

	workspace := occson.Workspace{Token: token}

	decrypted := workspace.Download(url, passphrase)
	fmt.Println(string(decrypted))

	blob := `[config]

param = "some_param"`

	workspace.Upload(url, passphrase, blob, true)
}
```
