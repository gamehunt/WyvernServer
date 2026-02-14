package domain

import (
	"wyvern/server/internal/types"
)

type Role struct {
	Id            types.ID       `bson:"id"`
	GuildId       types.ID       `bson:"guild_id"`
	Name          string         `bson:"name"`
	Color         [3]int         `bson:"color"`
	Permissions   Permission     `bson:"permissions"`
	Position      int            `bson:"position"`
}
