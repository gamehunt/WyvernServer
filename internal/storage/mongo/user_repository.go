package mongo

import (
	"context"
	"time"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UsersRepositoryImpl struct {
    collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, name string) repository.UserRepository {
    return &UsersRepositoryImpl{
        collection: client.Database(name).Collection("users"),
    }
}

func (r *UsersRepositoryImpl) Create(msg *domain.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, msg)
    return err
}

func (r *UsersRepositoryImpl) FindByIdentity(identity []byte) (*domain.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var user domain.User
    err := r.collection.FindOne(ctx, bson.M{"opaque_record.client_identity": identity}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &user, nil
}
