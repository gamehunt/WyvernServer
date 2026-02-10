package domain

import (
    "time"
    "github.com/bytemare/opaque"
)

type User struct {
	OpaqueRecord  *opaque.ClientRecord
	DisplayName   string
	JoinedAt      time.Time
	LastOnlineAt  time.Time
}
