package domain

import (
	"time"

	"github.com/bytemare/opaque"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	OpaqueRecord  *opaque.ClientRecord `bson:"opaque_record,omitempty"`
	DisplayName   string               `bson:"display_name"`
	JoinedAt      time.Time            `bson:"joined_at"`
	LastOnlineAt  time.Time            `bson:"last_online_at"`
}
