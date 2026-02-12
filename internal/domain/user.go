package domain

import (
	"time"
	"wyvern/server/internal/types"
)

type User struct {
	Id            types.ID  `bson:"id"`
	Identity      string    `bson:"identity"`
	DisplayName   string    `bson:"display_name,omitempty"`
	OpaqueRecord  *[]byte   `bson:"opaque_record,omitempty"`
	JoinedAt      time.Time `bson:"joined_at"`
	LastOnlineAt  time.Time `bson:"last_online_at"`
}
