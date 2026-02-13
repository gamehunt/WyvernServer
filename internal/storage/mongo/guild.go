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

type GuildRepositoryImpl struct {
    collection *mongo.Collection
}

func NewGuildRepository(client *mongo.Client, name string) repository.GuildRepository {
    return &GuildRepositoryImpl{
        collection: client.Database(name).Collection("guilds"),
    }
}

func (r *GuildRepositoryImpl) Create(guild *domain.Guild) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, guild)
    return err
}

func (r *GuildRepositoryImpl) Update(guild *domain.Guild) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	filter := bson.D{{Key: "id", Value: guild.Id}}

	_, err := r.collection.ReplaceOne(ctx, filter, guild)

    return err
}

func (r *GuildRepositoryImpl) Delete(id types.ID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *GuildRepositoryImpl) FindByID(id types.ID) (*domain.Guild, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var guild domain.Guild
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&guild)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &guild, nil
}
