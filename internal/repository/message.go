package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type MessageRepository interface {
	Create(channel *domain.Message) error
	Update(channel *domain.Message) error
	Delete(id types.ID)             error

	FindByID(id types.ID)           (*domain.Message,  error)
	FindByChannel(chanId types.ID)  ([]domain.Message, error)
}
