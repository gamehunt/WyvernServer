package service

import (
	"time"
	"wyvern/server/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {

}

func NewTokenService() (*TokenService, error) {
	return &TokenService{}, nil
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

func (s *TokenService) VerifyAccessToken(input string, sessionKey []byte) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(input, claims, func(token *jwt.Token) (any, error) {
    	return sessionKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	return token, claims, nil
}
