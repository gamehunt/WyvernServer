package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type OpaqueSessionCache interface {
	Save(*domain.OpaqueLoginSession) (*types.ID, error)
	Restore(types.ID) (*domain.OpaqueLoginSession, error)
}
