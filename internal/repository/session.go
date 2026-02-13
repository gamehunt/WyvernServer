package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type SessionRepository interface {
	Create(session *domain.Session) error
	Update(session *domain.Session) error
	Delete(id types.ID) error

	FindByID(id types.ID) (*domain.Session, error)
}
