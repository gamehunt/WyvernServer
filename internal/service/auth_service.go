package service

import (
	"time"
	"wyvern/server/internal/config"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/pkg/errors"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/storage"
	"wyvern/server/internal/types"

	"github.com/bytemare/opaque"
)

type AuthService struct {
    usersRepo repository.UserRepository
	opaqueSessionCache repository.OpaqueSessionCache

	identity  []byte
	secrets   domain.Secret
	conf     *opaque.Configuration
}

func NewAuthService(
    usersRepo repository.UserRepository,
	opaqueSessionCache repository.OpaqueSessionCache,	
	authCfg config.AuthConfig,
) *AuthService {
	conf := opaque.DefaultConfiguration()
	secrets := storage.GetSecrets(authCfg.SecretPath, conf)
    return &AuthService{
        usersRepo: usersRepo,
		conf: opaque.DefaultConfiguration(),
		secrets: secrets,
		identity: []byte(authCfg.ServerIdentity),
		opaqueSessionCache: opaqueSessionCache,
    }
}

func (s *AuthService) StartLogin(identity string, msg1 []byte) ([]byte, []byte, *types.ID, error) {
	server, err := s.conf.Server()
	if err != nil {
		return nil, nil, nil, err
	}

	err = server.SetKeyMaterial(s.identity, s.secrets.ServerPrivateKey, s.secrets.ServerPublicKey, s.secrets.SecretOprfSeed)
	if err != nil {
		return nil, nil, nil, err
	}

	ke1, err := server.Deserialize.KE1(msg1)
	if err != nil {
		return nil, nil, nil, err
	}

	clientRecord, err := s.usersRepo.FindByIdentity(identity)
	if clientRecord == nil {
		return nil, nil, nil, errors.AuthInvalidIdentity
	}

	record, err := server.Deserialize.RegistrationRecord(*clientRecord.OpaqueRecord)
	if err != nil {
		return nil, nil, nil, err
	}

	opaqueClientRecord := &opaque.ClientRecord{
		CredentialIdentifier: []byte(clientRecord.Id.String()),
		ClientIdentity:       []byte(clientRecord.Identity),
		RegistrationRecord:   record,
	}

	ke2, err := server.LoginInit(ke1, opaqueClientRecord)
	if err != nil {
		return nil, nil, nil, err
	}

	opaqueSession := &domain.OpaqueLoginSession{
		UserId: clientRecord.Id,
		Data: server.SerializeState(),
	}

	sessionId, err := s.opaqueSessionCache.Save(opaqueSession)	
	if err != nil {
		return nil, nil, nil, err
	}

	return ke2.Serialize(), s.identity, sessionId, nil
}

func (s *AuthService) FinishLogin(opaqueSessionId types.ID, msg3 []byte) ([]byte, *types.ID, error) {
	server, err := s.conf.Server()

	if err != nil {
		return nil, nil, err
	}

	err = server.SetKeyMaterial(s.identity, s.secrets.ServerPrivateKey, s.secrets.ServerPublicKey, s.secrets.SecretOprfSeed)
	if err != nil {
		return nil, nil, err
	}

	opaqueSession, err := s.opaqueSessionCache.Restore(opaqueSessionId)
	if err != nil {
		return nil, nil, err
	}

	server.SetAKEState(opaqueSession.Data)

	ke3, err := server.Deserialize.KE3(msg3)
	if err != nil {
		return nil, nil, err
	}

	if err := server.LoginFinish(ke3); err != nil {
		return nil, nil, err
	}

	return server.SessionKey(), &opaqueSession.UserId, nil
}

func (s *AuthService) StartRegister(msg1 []byte) ([]byte, []byte, *types.ID, error) {
	server, err := s.conf.Server()
	if err != nil {
		return nil, nil, nil, err
	}

	err = server.SetKeyMaterial(s.identity, s.secrets.ServerPrivateKey, s.secrets.ServerPublicKey, s.secrets.SecretOprfSeed)
	if err != nil {
		return nil, nil, nil, err
	}

	request, err := server.Deserialize.RegistrationRequest(msg1)
	if err != nil {
		return nil, nil, nil, err
	}

	credID, err := types.NewID()
	if err != nil {
		return nil, nil, nil, err
	}

	user := domain.User{
		Id: credID,
	}

	err = s.usersRepo.Create(&user)
	if err != nil {
		return nil, nil, nil, err
	}

	pks, err := server.Deserialize.DecodeAkePublicKey(s.secrets.ServerPublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	response := server.RegistrationResponse(request, pks, []byte(credID.String()), s.secrets.SecretOprfSeed)

	return response.Serialize(), s.identity, &credID, nil
}

func (s *AuthService) FinishRegister(credID types.ID, identity string, msg3 []byte) error {
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

	user, err := s.usersRepo.FindByID(credID)
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
