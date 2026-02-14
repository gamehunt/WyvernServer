package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type RoleRepository interface {
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	Delete(id types.ID) error

	FindByID(id types.ID) (*domain.Role, error)
	FindByIDs(id []types.ID) ([]domain.Role, error)
	FindByGuild(guildID types.ID) ([]domain.Role, error)
}
