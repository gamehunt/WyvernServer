package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type SessionCache interface {
	Save(session *domain.Session)  error
	Get(sessionId types.ID)        (*domain.Session, error)
	Invalidate(sessionId types.ID) error
}
