package domain

import "wyvern/server/internal/types"

type Guild struct {
	Id      types.ID `bson:"id"`
	Name    string   `bson:"name"`
	OwnerId types.ID `bson:"owner_id"`
}
