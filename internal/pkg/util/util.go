package util

import (
	"encoding/base64"
	"fmt"

	"github.com/bytemare/opaque"
	"github.com/gin-gonic/gin"
)

func RandomBase64(bytes int) string {
	return base64.RawURLEncoding.EncodeToString(opaque.RandomBytes(bytes))
}

func HttpError(ctx *gin.Context, code int, err error) {
    ctx.JSON(code, gin.H{"error": err.Error()})
}

func HttpErrorString(ctx *gin.Context, code int, err string) {
	HttpError(ctx, code, fmt.Errorf(err))
}
