package occson

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
