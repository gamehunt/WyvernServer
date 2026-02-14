package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type GuildRepository interface {
	Create(guild *domain.Guild) error
	Update(guild *domain.Guild) error
	Delete(id types.ID)         error

	FindByID(id types.ID) (*domain.Guild, error)
	FindByIDs(ids []types.ID) ([]domain.Guild, error)
}
