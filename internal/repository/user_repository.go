package repository

import (
	domain "wyvern/server/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
	FindByID(id bson.ObjectID) (*domain.User, error)
	FindByIdentity(identity string) (*domain.User, error)
}
