package redis

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"

	"github.com/valkey-io/valkey-go"
)

type KeyringStorageImpl struct {
	client valkey.Client
}

func NewKeyringStorage(client valkey.Client) repository.KeyringStorage {
    return &KeyringStorageImpl{
		client: client,
    }
}

func (s *KeyringStorageImpl) Key(sessionId types.ID) string {
	return "keyring:" + sessionId.String()
}

func (s *KeyringStorageImpl) Get(sessionId types.ID) (*domain.Keyring, error) {
	cmd := s.client.B().JsonGet().Key(s.Key(sessionId)).Path("$").Build()

	result := s.client.Do(context.Background(), cmd)

	if result.Error() != nil {
		return nil, result.Error()
	}

	var keyring []domain.Keyring

	err := result.DecodeJSON(&keyring)
	if err != nil {
		return nil, err
	}

	if len(keyring) != 1 {
		return nil, fmt.Errorf("Invalid keyring array")
	}

	return &keyring[0], nil
}

func (s *KeyringStorageImpl) Delete(sessionId types.ID) error {
	cmd := s.client.B().JsonDel().Key(s.Key(sessionId)).Path("$").Build()
	return s.client.Do(context.Background(), cmd).Error()
}

func (s *KeyringStorageImpl) derive(sessionKey []byte) ([]byte, []byte) {
	h := sha256.New()
    h.Write(sessionKey)
    h.Write([]byte("auth"))
	authKey := h.Sum(nil)

    h.Reset()
    h.Write(sessionKey)
    h.Write([]byte("sign"))
	signKey := h.Sum(nil)

	return authKey, signKey
}

func (s *KeyringStorageImpl) Generate(sessionId types.ID, sessionKey []byte, ttl time.Duration) (*domain.Keyring, error) {
	ak, sk := s.derive(sessionKey)
	keyring := &domain.Keyring{
		AuthKey: ak,
		SignKey: sk,
	}

	jsonData, err := json.Marshal(keyring)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	setCmd := s.client.B().
	JsonSet().
		Key(s.Key(sessionId)).
		Path("$").
		Value(string(jsonData)).Build()

	err = s.client.Do(ctx, setCmd).Error()

	if err != nil {
		return nil, err
	}

	err = s.Refresh(sessionId, ttl)

	if err != nil {
		return nil, err
	}

	return keyring, nil
}

func (s *KeyringStorageImpl) Refresh(sessionId types.ID, ttl time.Duration) error {
	expireCmd := s.client.B().Expire().Key(s.Key(sessionId)).Seconds(int64(ttl.Seconds())).Build()

	err := s.client.Do(context.Background(), expireCmd).Error()

	if err != nil {
		return err
	}

	return nil
}
