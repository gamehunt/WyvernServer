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

func NewOpaqueSessionStorage(client valkey.Client) *OpaqueSessionStorage {
    return &OpaqueSessionStorage{
		client: client,
    }
}

func (r *OpaqueSessionStorage) Save(state []byte) (string, error) {
	sessionId := base64.RawURLEncoding.EncodeToString(opaque.RandomBytes(32))

	ctx := context.Background()

	err := r.client.Do(ctx, 
	r.client.B().Set().Key("opaque:"+sessionId).Value(base64.RawURLEncoding.EncodeToString(state)).Nx().Build()).Error()

	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (r *OpaqueSessionStorage) Restore(sessionId string) ([]byte, error) {
	result := r.client.Do(context.Background(), r.client.B().Get().Key("opaque:"+sessionId).Build())

	if (result.Error() != nil) {
		return nil, result.Error()
	}

	encState, err := result.AsBytes()
	if err != nil {
		return nil, err
	}

	state, err := base64.RawURLEncoding.DecodeString(string(encState))
	if err != nil {
		return nil, err
	}

	return state, nil
}
