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

type MessageRepositoryImpl struct {
    collection *mongo.Collection
}

func NewMessageRepository(client *mongo.Client, name string) repository.MessageRepository {
    return &MessageRepositoryImpl{
        collection: client.Database(name).Collection("channels"),
    }
}

func (r *MessageRepositoryImpl) Create(message *domain.Message) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, message)
    return err
}

func (r *MessageRepositoryImpl) Update(message *domain.Message) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	filter := bson.D{{Key: "id", Value: message.Id}}

	_, err := r.collection.ReplaceOne(ctx, filter, message)

    return err
}

func (r *MessageRepositoryImpl) Delete(id types.ID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *MessageRepositoryImpl) FindByID(id types.ID) (*domain.Message, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var message domain.Message
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&message)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &message, nil
}

func (r *MessageRepositoryImpl)	FindByChannel(chanId types.ID) ([]domain.Message, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	entries, err := r.collection.Find(ctx, bson.M{"channel_id": chanId})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	defer entries.Close(ctx)

	var messages []domain.Message
    for entries.Next(context.TODO()) {
        var msg domain.Message
        err := entries.Decode(&msg)
        if err != nil {
			return nil, err
        }

        messages = append(messages, msg)
    }

    return messages, nil
}
