package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type PermissionOverrideRepository interface {
	Create(permOverride *domain.PermissionOverride) error
	Update(permOverride *domain.PermissionOverride) error
	Delete(id types.ID) error

	FindByID(id types.ID) (*domain.PermissionOverride, error)
	FindByChannel(channel types.ID) ([]domain.PermissionOverride, error)
}
