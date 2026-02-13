package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type ChannelRepository interface {
	Create(channel *domain.Channel) error
	Update(channel *domain.Channel) error
	Delete(id types.ID)             error

	FindByID(id types.ID)           (*domain.Channel,  error)
	FindByGuild(guildId types.ID)   ([]domain.Channel, error)
}
