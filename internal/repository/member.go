package repository

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/types"
)

type MemberRepository interface {
	Create(member *domain.Member) error
	Update(member *domain.Member) error
	Delete(userId types.ID, guildId types.ID) error

	FindByGuild(guildId types.ID)  ([]domain.Member, error)
	GetForUser(userId types.ID, guildId types.ID) (*domain.Member, error)
	GetUserGuilds(userId types.ID) ([]domain.Guild, error)
}
