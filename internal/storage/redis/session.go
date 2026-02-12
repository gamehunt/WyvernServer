package redis

import (
	"context"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"

	"github.com/valkey-io/valkey-go"
)

type SessionCacheImpl struct {
	client valkey.Client
}

func NewSessionCache(client valkey.Client) repository.SessionCache {
    return &SessionCacheImpl{
		client: client,
    }
}

func (r *SessionCacheImpl) Save(session *domain.Session) error {
	return r.client.Do(context.Background(), r.client.
	B().
	Set().
		Key("session:" + session.Id.String()).
		Value(session.UserId.String()).
	Build()).
	Error()
}

func (r *SessionCacheImpl) Get(sessionId types.ID) (*domain.Session, error) {
	result := r.client.Do(context.Background(), r.client.
	B().
	Get().
		Key("session:"+sessionId.String()).
	Build())

	userIdBytes, err := result.AsBytes()
	if err != nil {
		return nil, err
	}

	receivedUserId, err := types.ParseID(string(userIdBytes))
	if err != nil {
		return nil, err
	}

	return &domain.Session{
		Id:     sessionId,
		UserId: receivedUserId,
	}, nil
}

func (r *SessionCacheImpl) Invalidate(sessionId types.ID) error {
	return r.client.Do(context.Background(), r.client.B().Del().Key("session:"+sessionId.String()).Build()).Error()
}
