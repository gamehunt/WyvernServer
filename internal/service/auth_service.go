package service

import (
	"time"
	"wyvern/server/internal/config"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/pkg/errors"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/storage"

	"github.com/bytemare/opaque"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthService struct {
    usersRepo repository.UserRepository

	identity  []byte
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
		identity: []byte(authCfg.ServerIdentity),
    }
}

func (s *AuthService) StartLogin(identity string, msg1 []byte) ([]byte, []byte, error) {
	server, err := s.conf.Server()
	if err != nil {
		return nil, nil, err
	}

	err = server.SetKeyMaterial(s.identity, s.secrets.ServerPrivateKey, s.secrets.ServerPublicKey, s.secrets.SecretOprfSeed)
	if err != nil {
		return nil, nil, err
	}

	ke1, err := server.Deserialize.KE1(msg1)
	if err != nil {
		return nil, nil, err
	}

	clientRecord, err := s.usersRepo.FindByIdentity(identity)
	if clientRecord == nil {
		return nil, nil, errors.AuthInvalidIdentity
	}

	record, err := server.Deserialize.RegistrationRecord(*clientRecord.OpaqueRecord)
	if err != nil {
		return nil, nil, err
	}

	opaqueClientRecord := &opaque.ClientRecord{
		CredentialIdentifier: []byte(clientRecord.ID.Hex()),
		ClientIdentity:       []byte(clientRecord.Identity),
		RegistrationRecord:   record,
	}

	ke2, err := server.LoginInit(ke1, opaqueClientRecord)
	if err != nil {
		return nil, nil, err
	}

	return ke2.Serialize(), s.identity, nil
}

func (s *AuthService) FinishLogin(msg3 []byte) ([]byte, error) {
	server, err := s.conf.Server()

	if err != nil {
		return nil, err
	}

	err = server.SetKeyMaterial(s.identity, s.secrets.ServerPrivateKey, s.secrets.ServerPublicKey, s.secrets.SecretOprfSeed)
	if err != nil {
		return nil, err
	}

	ke3, err := server.Deserialize.KE3(msg3)
	if err != nil {
		return nil, err
	}

	if err := server.LoginFinish(ke3); err != nil {
		return nil, err
	}

	return server.SessionKey(), nil
}

func (s *AuthService) StartRegister(msg1 []byte) ([]byte, []byte, string, error) {
	server, err := s.conf.Server()
	if err != nil {
		return nil, nil, "", err
	}

	err = server.SetKeyMaterial(s.identity, s.secrets.ServerPrivateKey, s.secrets.ServerPublicKey, s.secrets.SecretOprfSeed)
	if err != nil {
		return nil, nil, "", err
	}

	request, err := server.Deserialize.RegistrationRequest(msg1)
	if err != nil {
		return nil, nil, "", err
	}

	user := domain.User{
		ID: bson.NewObjectID(),
	}

	credID := user.ID.Hex()

	err = s.usersRepo.Create(&user)
	if err != nil {
		return nil, nil, "", err
	}

	pks, err := server.Deserialize.DecodeAkePublicKey(s.secrets.ServerPublicKey)
	if err != nil {
		return nil, nil, "", err
	}

	response := server.RegistrationResponse(request, pks, []byte(credID), s.secrets.SecretOprfSeed)

	return response.Serialize(), s.identity, credID, nil
}

func (s *AuthService) FinishRegister(credID string, identity string, msg3 []byte) error {
	server, err := s.conf.Server()
	if err != nil {
		return err
	}

	err = server.SetKeyMaterial(s.identity, s.secrets.ServerPrivateKey, s.secrets.ServerPublicKey, s.secrets.SecretOprfSeed)
	if err != nil {
		return err
	}

	record, err := server.Deserialize.RegistrationRecord(msg3)
	if err != nil {
		return err
	}

	idObj, err := bson.ObjectIDFromHex(credID)
	if err != nil {
		return err
	}

	user, err := s.usersRepo.FindByID(idObj)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.AuthInvalidCredId
	}

	existingUser, err := s.usersRepo.FindByIdentity(identity)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.UserAlreadyExists
	}

	opaqueRecordBytes := record.Serialize()

	user.Identity     = identity
	user.OpaqueRecord = &opaqueRecordBytes
	user.JoinedAt     = time.Now()
	user.LastOnlineAt = user.JoinedAt

	return s.usersRepo.Update(user)
}
