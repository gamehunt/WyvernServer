package types

import (
	"crypto/sha256"
	"time"
	"wyvern/server/internal/pkg/util"
)

type RefreshToken struct {
	Hash      []byte    `bson:"hash"`
	ExpiresAt time.Time `bson:"expires_at"`
}

func NewRefreshToken(lifetime time.Duration) (string, RefreshToken) {
	src := util.RandomBase64(32)

	hashBytes := sha256.Sum256([]byte(src))
	return src, RefreshToken {
		Hash: hashBytes[:],
		ExpiresAt: time.Now().Add(lifetime),
	}
}

func (tok RefreshToken) Verify(in string) bool {
	if tok.ExpiresAt.Before(tok.ExpiresAt) {
		return false
	}

	inHash := sha256.Sum256([]byte(in))

	return inHash == [32]byte(tok.Hash)
}
