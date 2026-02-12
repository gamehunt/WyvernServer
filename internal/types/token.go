package types

import (
	"crypto/sha256"
	"time"
	"wyvern/server/internal/pkg/util"
)

type RefreshToken []byte

func NewRefreshToken(lifetime time.Duration) (string, RefreshToken) {
	src := util.RandomBase64(32)

	hashBytes := sha256.Sum256([]byte(src))
	return src, hashBytes[:]
}

func (tok RefreshToken) Verify(in string) bool {
	inHash := sha256.Sum256([]byte(in))
	return inHash == [32]byte(tok)
}
