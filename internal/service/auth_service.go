package service

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/storage"

	"github.com/bytemare/opaque"
)

type AuthService struct {
    usersRepo repository.UserRepository

	secrets domain.Secret
	conf *opaque.Configuration
}

func NewAuthService(
    usersRepo repository.UserRepository,
) *AuthService {
	conf := opaque.DefaultConfiguration()
	secrets := storage.GetSecrets(conf)
    return &AuthService{
        usersRepo: usersRepo,
		conf: opaque.DefaultConfiguration(),
		secrets: secrets,
    }
}

func (s *AuthService) Register(login string, password string) error {
   return nil 
}

func (s *AuthService) Login(login string, password string) error {
   return nil 
}
