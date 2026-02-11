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

func (r *UsersRepositoryImpl) Update(user *domain.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	filter := bson.D{{Key: "_id", Value: user.ID}}

	_, err := r.collection.ReplaceOne(ctx, filter, user)

    return err
}

func (r *UsersRepositoryImpl) FindByID(id bson.ObjectID) (*domain.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var user domain.User
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &user, nil
}

func (r *UsersRepositoryImpl) FindByIdentity(id string) (*domain.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var user domain.User
    err := r.collection.FindOne(ctx, bson.M{"identity": id}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &user, nil
}
