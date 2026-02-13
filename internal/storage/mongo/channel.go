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

type ChannelRepositoryImpl struct {
    collection *mongo.Collection
}

func NewChannelRepository(client *mongo.Client, name string) repository.ChannelRepository {
    return &ChannelRepositoryImpl{
        collection: client.Database(name).Collection("channels"),
    }
}

func (r *ChannelRepositoryImpl) Create(channel *domain.Channel) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, channel)
    return err
}

func (r *ChannelRepositoryImpl) Update(channel *domain.Channel) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	filter := bson.D{{Key: "id", Value: channel.Id}}

	_, err := r.collection.ReplaceOne(ctx, filter, channel)

    return err
}

func (r *ChannelRepositoryImpl) Delete(id types.ID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *ChannelRepositoryImpl) FindByID(id types.ID) (*domain.Channel, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var channel domain.Channel
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&channel)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &channel, nil
}

func (r *ChannelRepositoryImpl)	FindByGuild(guildId types.ID) ([]domain.Channel, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	entries, err := r.collection.Find(ctx, bson.M{"guild_id": guildId})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var channels []domain.Channel
    for entries.Next(context.TODO()) {
        var ch domain.Channel
        err := entries.Decode(&ch)
        if err != nil {
			return nil, err
        }

        channels = append(channels, ch)
    }

    return channels, nil
}
