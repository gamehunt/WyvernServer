package domain

type Secret struct {
	SecretOprfSeed   []byte `json:"secret_oprf_seed"`
	ServerPrivateKey []byte `json:"server_private_key"`
	ServerPublicKey  []byte `json:"server_public_key"`
}
