package mongo

import (
	"context"
	"time"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UsersRepositoryImpl struct {
    collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) repository.UserRepository {
    return &UsersRepositoryImpl{
        collection: db.Collection("users"),
    }
}

func (r *UsersRepositoryImpl) Create(msg *domain.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, msg)
    return err
}
