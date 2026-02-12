package repository

import (
	"time"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type KeyringStorage interface {
	Get(types.ID) (*domain.Keyring, error)
	Delete(types.ID) error
	Generate(types.ID, []byte, time.Duration) (*domain.Keyring, error)
	Refresh(types.ID, time.Duration) error
}
