package mongo

import (
	"context"
	"time"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


type PermissionOverrideRepositoryImpl struct {
    collection *mongo.Collection
}

func NewPermissionOverrideRepository(client *mongo.Client, name string) repository.PermissionOverrideRepository {
    return &PermissionOverrideRepositoryImpl{
        collection: client.Database(name).Collection("permission_overrides"),
    }
}

func (r *PermissionOverrideRepositoryImpl) Create(guild *domain.PermissionOverride) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, guild)
    return err
}

func (r *PermissionOverrideRepositoryImpl) Update(guild *domain.PermissionOverride) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	filter := bson.D{{Key: "id", Value: guild.Id}}

	_, err := r.collection.ReplaceOne(ctx, filter, guild)

    return err
}

func (r *PermissionOverrideRepositoryImpl) Delete(id types.ID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *PermissionOverrideRepositoryImpl) FindByID(id types.ID) (*domain.PermissionOverride, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var override domain.PermissionOverride
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&override)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &override, nil
}

func (r *PermissionOverrideRepositoryImpl) FindByChannel(channelId types.ID) ([]domain.PermissionOverride, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	entries, err := r.collection.Find(ctx, bson.M{"channel_id": channelId})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	defer entries.Close(ctx)

	var permOverrides []domain.PermissionOverride
    for entries.Next(context.TODO()) {
        var r domain.PermissionOverride
        err := entries.Decode(&r)
        if err != nil {
			return nil, err
        }

        permOverrides = append(permOverrides, r)
    }

    return permOverrides, nil
}
