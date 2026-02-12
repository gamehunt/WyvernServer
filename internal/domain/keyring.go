package domain

type Keyring struct {
	AuthKey []byte `json:"auth_key"`
	SignKey []byte `json:"sign_key"`
}
