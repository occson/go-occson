# go-occson

```go
func main() {
	url := "ccs://golang-test.toml"
	token := "decafc0ffeebad"
	passphrase := "deadbeef"

	workspace := Workspace{Token: token}

	decrypted := workspace.Download(url, passphrase)
	fmt.Println(string(decrypted))

	blob := `[config]

param = "yurp"`

	workspace.Upload(url, passphrase, blob, true)
}
```
