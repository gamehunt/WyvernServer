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

type SessionRepositoryImpl struct {
    collection *mongo.Collection
}

func NewSessionRepository(client *mongo.Client, name string) repository.SessionRepository {
    return &SessionRepositoryImpl{
        collection: client.Database(name).Collection("sessions"),
    }
}

func (r *SessionRepositoryImpl) Create(session *domain.Session) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, session)
    return err
}

func (r *SessionRepositoryImpl) Update(session *domain.Session) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	filter := bson.D{{Key: "id", Value: session.Id}}

	_, err := r.collection.ReplaceOne(ctx, filter, session)

    return err
}

func (r *SessionRepositoryImpl) FindByID(id types.ID) (*domain.Session, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var session domain.Session
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&session)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &session, nil
}

func (r *SessionRepositoryImpl) Delete(id types.ID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

