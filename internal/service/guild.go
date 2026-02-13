package service

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"
)

type GuildService struct {
	guildRepo repository.GuildRepository
}

func NewGuildService(
	guildRepo repository.GuildRepository,
) *GuildService {
	return &GuildService {
		guildRepo: guildRepo,
	}	
}

func (s *GuildService) GetGuild(id types.ID) (*domain.Guild, error) {
	return s.guildRepo.FindByID(id)
}
