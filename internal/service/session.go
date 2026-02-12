package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"time"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"
)

type SessionService struct {
	keyring      repository.KeyringStorage
	persistent   repository.SessionRepository
	tokenService *TokenService
}

func NewSessionService(
	keyring        repository.KeyringStorage,
	sessionStorage repository.SessionRepository,
	tokenService   *TokenService,
) *SessionService {
	return &SessionService{
		keyring: keyring,
		persistent: sessionStorage,
		tokenService: tokenService,
	}
}

func (r *SessionService) NewSession(userId types.ID, sessionKey []byte) (*domain.Session, *string, *string, error) {
	id, err := types.NewID() 
	if err != nil {
		return nil, nil, nil, err
	}

	ttl := 30 * 24 * time.Hour

	refreshTokStr, refreshTok := types.NewRefreshToken(ttl)

	session := &domain.Session{
		Id: id,
		UserId: userId,
		RefreshToken: refreshTok,
		ExpiresAt: time.Now().Add(ttl),
		LastActive: time.Now(),
	}

	err = r.persistent.Create(session)
	if err != nil {
		return nil, nil, nil, err
	}

	keys, err := r.keyring.Generate(session.Id, sessionKey, ttl)
	if err != nil {
		return nil, nil, nil, err
	}

	accessToken, err := r.tokenService.CreateAccessToken(session, keys.SignKey, 15 * time.Minute)
	if err != nil {
		return nil, nil, nil, err
	}

	return session, &refreshTokStr, &accessToken, nil
}

func (r *SessionService) RefreshSession(sessionId types.ID, refreshToken string, proof []byte, nonce []byte) (*string, error) {
	keys, err := r.keyring.Get(sessionId)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
    h.Write(keys.AuthKey)
    h.Write([]byte(nonce))

	expected := h.Sum(nil)
	
	if !hmac.Equal(expected, proof) {
		return nil, fmt.Errorf("Invalid proof")
	}

	session, err := r.persistent.FindByID(sessionId)
	if err != nil {
		return nil, err
	}

	if !session.RefreshToken.Verify(refreshToken) {
		return nil, fmt.Errorf("Invalid refresh token")
	}

	ttl := 15 * time.Minute

	err = r.keyring.Refresh(session.Id, ttl)
	if err != nil {
		return nil, err
	}

	accessToken, err := r.tokenService.CreateAccessToken(session, keys.SignKey, ttl)
	if err != nil {
		return nil, err
	}

	ttl = 30 * 24 * time.Hour

	session.ExpiresAt = time.Now().Add(ttl)
	err = r.persistent.Update(session)

	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}
