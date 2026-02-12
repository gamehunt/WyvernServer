package middleware

import (
	"net/http"
	"strings"
	"wyvern/server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(tokenService *service.TokenService) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

		words := strings.Split(header, " ")

		if len(words) != 2 {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
		}

		if words[0] != "Bearer" {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
		}

		token := words[1]

		jwtToken, err := tokenService.VerifyAccessToken(token)
		if err != nil || !jwtToken.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
        sessionID, _ := claims["sid"].(string)
		userID, _ := claims["sub"].(string)


        c.Set("sessionId", sessionID)
        c.Set("userId", userID)
        c.Next()
    }
}
