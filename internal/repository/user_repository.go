package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type UserRepository interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
	FindByID(id types.ID) (*domain.User, error)
	FindByIdentity(identity string) (*domain.User, error)
}
