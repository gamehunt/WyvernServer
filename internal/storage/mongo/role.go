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

type RoleRepositoryImpl struct {
    collection *mongo.Collection
}

func NewRoleRepository(client *mongo.Client, name string) repository.RoleRepository {
    return &RoleRepositoryImpl{
        collection: client.Database(name).Collection("roles"),
    }
}

func (r *RoleRepositoryImpl) Create(guild *domain.Role) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    _, err := r.collection.InsertOne(ctx, guild)
    return err
}

func (r *RoleRepositoryImpl) Update(guild *domain.Role) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	filter := bson.D{{Key: "id", Value: guild.Id}}

	_, err := r.collection.ReplaceOne(ctx, filter, guild)

    return err
}

func (r *RoleRepositoryImpl) Delete(id types.ID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *RoleRepositoryImpl) FindByID(id types.ID) (*domain.Role, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	var role domain.Role
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&role)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

    return &role, nil
}


func (r *RoleRepositoryImpl) FindByIDs(ids []types.ID) ([]domain.Role, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	entries, err := r.collection.Find(ctx, bson.M{"id": bson.M{"$in": ids}})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	defer entries.Close(ctx)

	var roles []domain.Role
    for entries.Next(context.TODO()) {
        var r domain.Role
        err := entries.Decode(&r)
        if err != nil {
			return nil, err
        }

        roles = append(roles, r)
    }

    return roles, nil
}

func (r *RoleRepositoryImpl) FindByGuild(guildId types.ID) ([]domain.Role, error) {
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

	var roles []domain.Role
    for entries.Next(context.TODO()) {
        var r domain.Role
        err := entries.Decode(&r)
        if err != nil {
			return nil, err
        }

        roles = append(roles, r)
    }

    return roles, nil
}
