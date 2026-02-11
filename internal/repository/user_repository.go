package repository

import (
	domain "wyvern/server/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByIdentity(identity []byte) (*domain.User, error)
}
