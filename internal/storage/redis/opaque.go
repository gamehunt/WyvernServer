package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"

	"github.com/valkey-io/valkey-go"
)

type OpaqueSessionCacheImpl struct {
	client valkey.Client
}

func NewOpaqueSessionCache(client valkey.Client) repository.OpaqueSessionCache {
    return &OpaqueSessionCacheImpl{
		client: client,
    }
}

func (r *OpaqueSessionCacheImpl) Save(session *domain.OpaqueLoginSession) (*types.ID, error) {
	sessionId, err := types.NewID()
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	err = r.client.Do(ctx, 
	r.client.B().
	JsonSet().
		Key("opaque:"+sessionId.String()).
		Path("$").
		Value(string(jsonData)).
	Nx().Build()).Error()

	if err != nil {
		return nil, err
	}

	return &sessionId, nil
}

func (r *OpaqueSessionCacheImpl) Restore(sessionId types.ID) (*domain.OpaqueLoginSession, error) {
	sessionKey := "opaque:" + sessionId.String()

	result := r.client.Do(context.Background(),
	r.client.B().
		JsonGet().
			Key(sessionKey).
			Path("$").
		Build())

	if (result.Error() != nil) {
		return nil, result.Error()
	}

	var sessions []domain.OpaqueLoginSession

	err := result.DecodeJSON(&sessions)
	if err != nil {
		return nil, err
	}

	if len(sessions) != 1 {
		return nil, fmt.Errorf("Invalid sessions array received")
	}

	r.client.Do(context.Background(), r.client.B().JsonDel().Key(sessionKey).Path("$").Build())

	return &sessions[0], nil
}
