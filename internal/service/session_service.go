package service

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"
)

type SessionService struct {
	cache      repository.SessionCache
	persistent repository.SessionRepository
}

func NewSessionService(
	sessionCache   repository.SessionCache,
	sessionStorage repository.SessionRepository,
) *SessionService {
	return &SessionService{
		cache: sessionCache,
		persistent: sessionStorage,
	}
}

func (r *SessionService) NewSession(userId types.ID) (*domain.Session, error) {
	id, err := types.NewID() 
	if err != nil {
		return nil, err
	}

	session := &domain.Session{
		Id: id,
		UserId: userId,
	}

	err = r.persistent.Create(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionService) Get(sessionId types.ID) (*domain.Session, error) {
	session, err := r.cache.Get(sessionId)	
	if err == nil {
		return session, nil
	}

	session, err = r.persistent.FindByID(sessionId)
	if err != nil {
		return nil, err
	}

	r.cache.Save(session)

	return session, nil
}
