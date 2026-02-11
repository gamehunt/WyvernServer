package service

import (
	"wyvern/server/internal/config"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/pkg/errors"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/storage"

	"github.com/bytemare/opaque"
)

type AuthService struct {
    usersRepo repository.UserRepository

	identity  string
	secrets   domain.Secret
	conf     *opaque.Configuration
}

func NewAuthService(
    usersRepo repository.UserRepository,
	authCfg config.AuthConfig,
) *AuthService {
	conf := opaque.DefaultConfiguration()
	secrets := storage.GetSecrets(authCfg.SecretPath, conf)
    return &AuthService{
        usersRepo: usersRepo,
		conf: opaque.DefaultConfiguration(),
		secrets: secrets,
		identity: authCfg.ServerIdentity,
    }
}

func (s *AuthService) Register(login string, password string) error {
   return nil 
}


func (s *AuthService) StartLogin(identity []byte, msg1 []byte) ([]byte, error) {
	server, err := s.conf.Server()

	if err != nil {
		return nil, err
	}

	server.SetKeyMaterial([]byte(s.identity), s.secrets.SecretOprfSeed, s.secrets.ServerPrivateKey, s.secrets.ServerPublicKey)

	ke1, err := server.Deserialize.KE1(msg1)
	if err != nil {
		return nil, err
	}

	clientRecord, err := s.usersRepo.FindByIdentity(identity)

	if clientRecord == nil {
		return nil, errors.AuthInvalidIdentity
	}

	ke2, err := server.LoginInit(ke1, clientRecord.OpaqueRecord)
	if err != nil {
		return nil, err
	}

	return ke2.Serialize(), nil
}

func (s *AuthService) FinishLogin(msg3 []byte) (*[]byte, error) {
	server, err := s.conf.Server()

	if err != nil {
		return nil, err
	}

	server.SetKeyMaterial([]byte(s.identity), s.secrets.SecretOprfSeed, s.secrets.ServerPrivateKey, s.secrets.ServerPublicKey)

	ke3, err := server.Deserialize.KE3(msg3)
	if err != nil {
		return nil, err
	}

	if err := server.LoginFinish(ke3); err != nil {
		return nil, err
	}

	key := server.SessionKey()

	return &key, nil
}
