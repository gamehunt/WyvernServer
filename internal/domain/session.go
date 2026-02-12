package domain

import (
	"time"
	"wyvern/server/internal/types"
)

type Session struct {
	Id           types.ID           `bson:"id"`
	UserId       types.ID           `bson:"user_id"`
	RefreshToken types.RefreshToken `bson:"refresh"`
	LastActive   time.Time          `bson:"last_active"` 
	ExpiresAt    time.Time          `bson:"expires_at"`
}
