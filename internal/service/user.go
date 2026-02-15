package service

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"
)

type UserService struct {
	usersRepo repository.UserRepository
}

func NewUserService(usersRepo repository.UserRepository) *UserService {
	return &UserService{
		usersRepo: usersRepo,
	}
}

func (s *UserService) GetUser(userId types.ID) (*domain.User, error) {
	return s.usersRepo.FindByID(userId)
}
