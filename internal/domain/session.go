package domain

import "wyvern/server/internal/types"

type Session struct {
	Id     types.ID `bson:"id"`
	UserId types.ID `bson:"user_id"`
}
