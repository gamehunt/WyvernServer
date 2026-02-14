package domain

import "wyvern/server/internal/types"

type Guild struct {
	Id      types.ID `bson:"id" json:"id"`
	Name    string   `bson:"name" json:"name"`
	OwnerId types.ID `bson:"owner_id" json:"owner_id"`
}
