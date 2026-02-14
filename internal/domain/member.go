package domain

import (
	"time"
	"wyvern/server/internal/types"
)

type Member struct {
	UserId  types.ID   `bson:"user_id"`
	GuildId types.ID   `bson:"guild_id"`
	Nickname *string   `bson:"nickname,omitempty"`
	Roles   []types.ID `bson:"roles"`
	JoinedAt time.Time `bson:"joined_at"`
}
