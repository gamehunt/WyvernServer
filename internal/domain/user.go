package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID            bson.ObjectID `bson:"_id"`
	Identity      string        `bson:"identity"`
	DisplayName   string        `bson:"display_name,omitempty"`
	OpaqueRecord  *[]byte       `bson:"opaque_record,omitempty"`
	JoinedAt      time.Time     `bson:"joined_at"`
	LastOnlineAt  time.Time     `bson:"last_online_at"`
}
