package util

import (
	"encoding/base64"

	"github.com/bytemare/opaque"
)

func RandomBase64(bytes int) string {
	return base64.RawURLEncoding.EncodeToString(opaque.RandomBytes(bytes))
}
