package redis

import (
	"context"
	"encoding/base64"

	"github.com/bytemare/opaque"
	"github.com/valkey-io/valkey-go"
)

type OpaqueSessionStorage struct {
	client valkey.Client
}

func NewOpaqueSessionStorage(client valkey.Client, name string) *OpaqueSessionStorage {
    return &OpaqueSessionStorage{
		client: client,
    }
}

func (r *OpaqueSessionStorage) Save(state []byte) (string, error) {
	sessionId := base64.RawURLEncoding.EncodeToString(opaque.RandomBytes(32))

	ctx := context.Background()

  	// SET key val NX
	err := r.client.Do(ctx, r.client.B().Set().Key("opaque:"+sessionId).Value("val").Nx().Build()).Error()

	return "", nil
}
