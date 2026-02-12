package service

import (
	"time"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Keyring repository.KeyringStorage
}

func NewTokenService(keyring repository.KeyringStorage) *TokenService {
	return &TokenService{
		Keyring: keyring,
	}
}

func (s *TokenService) CreateAccessToken(session *domain.Session, sessionKey []byte, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": session.UserId.String(),	
		"sid": session.Id.String(),
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}	

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(sessionKey)
}

func (s *TokenService) VerifyAccessToken(input string) (*jwt.Token, error) {
	token, err := jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return nil, jwt.ErrTokenMalformed
        }

        sessionIDStr, ok := claims["sid"].(string)
        if !ok {
            return nil, jwt.ErrTokenMalformed
        }

		sessionID, err := types.ParseID(sessionIDStr)
		if err != nil {
			return nil, err
		}

        keys, err := s.Keyring.Get(sessionID)
        if err != nil {
            return nil, err
        }

        return keys.SignKey, nil
    })

	if err != nil {
		return nil, err
	}

	return token, nil
}
