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

type MemberRepositoryImpl struct {
    collection *mongo.Collection
}

func NewMemberRepository(client *mongo.Client, name string) repository.MemberRepository {
    return &MemberRepositoryImpl{
        collection: client.Database(name).Collection("members"),
    }
}

func (r *MemberRepositoryImpl) Create(member *domain.Member) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, member)
    return err
}

func (r *MemberRepositoryImpl) Update(member *domain.Member) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	filter := bson.D{{Key: "user_id", Value: member.UserId}, {Key: "guild_id", Value: member.GuildId}}

	_, err := r.collection.ReplaceOne(ctx, filter, member)

    return err
}

func (r *MemberRepositoryImpl) Delete(userId, guildId types.ID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.D{{Key: "user_id", Value: userId}, {Key: "guild_id", Value: guildId}})
	return err
}

func (r *MemberRepositoryImpl) FindByGuild(guildId types.ID) ([]domain.Member, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	entries, err := r.collection.Find(ctx, bson.M{"guild_id": guildId})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	defer entries.Close(ctx)

	var members []domain.Member
    for entries.Next(context.TODO()) {
        var ch domain.Member
        err := entries.Decode(&ch)
        if err != nil {
			return nil, err
        }

        members = append(members, ch)
    }

    return members, nil
}

func (r *MemberRepositoryImpl) GetForUser(userId, guildId types.ID) (*domain.Member, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var member domain.Member
	err := r.collection.FindOne(ctx, bson.D{{Key:"user_id", Value:userId}, 
	{Key: "guild_id", Value: guildId}}).Decode(&member)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &member, nil
}

func (r *MemberRepositoryImpl) FindByUser(userId types.ID) ([]domain.Member, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	entries, err := r.collection.Find(ctx, bson.M{"user_id": userId})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	defer entries.Close(ctx)

	var members []domain.Member
    for entries.Next(context.TODO()) {
        var ch domain.Member
        err := entries.Decode(&ch)
        if err != nil {
			return nil, err
        }

        members = append(members, ch)
    }

    return members, nil
}
