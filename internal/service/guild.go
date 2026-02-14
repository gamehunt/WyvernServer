package service

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/pkg/errors"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"
)

type GuildService struct {
	guildRepo  repository.GuildRepository
	membersSvc *MemberService
	permsSvc   *PermissionService
}

func NewGuildService(
	guildRepo  repository.GuildRepository,
	membersSvc *MemberService,
	permsSvc   *PermissionService,
) *GuildService {
	return &GuildService {
		guildRepo:  guildRepo,
		membersSvc: membersSvc,
		permsSvc:   permsSvc,
	}	
}

func (s *GuildService) GetGuild(id types.ID) (*domain.Guild, error) {
	return s.guildRepo.FindByID(id)
}

func (s *GuildService) CreateGuild(ownerId types.ID, name string) (*domain.Guild, error) {
	guildId, err := types.NewID()
	if err != nil {
		return nil, err
	}
	
	guild := &domain.Guild{
		Id: guildId,
		OwnerId: ownerId,
		Name: name,
	}

	err = s.guildRepo.Create(guild)
	if err != nil {
		return nil, err
	}

	_, err = s.membersSvc.JoinGuild(ownerId, guild.Id)
	if err != nil {
        return nil, err
	}

	return guild, nil
}

func (s *GuildService) UpdateGuild(requesterId, guildId types.ID, name string) error {
	guild, err := s.GetGuild(guildId)
	if err != nil {
		return err
	}

	if guild == nil {
		return errors.InvalidGuild
	}

	isAllowed, err := s.permsSvc.CanManageGuild(requesterId, guildId) 
	if err != nil {
		return err
	}

	if !isAllowed {
		return errors.InsufficientPerms
	}

	guild.Name = name

	err = s.guildRepo.Update(guild)
	if err != nil {
		return err
	}

	return nil
}

func (s *GuildService) DeleteGuild(requesterId, guildId types.ID) error {
	guild, err := s.guildRepo.FindByID(guildId)
    if err != nil {
        return err
    }

	if guild.OwnerId != requesterId {
		return errors.InsufficientPerms
	}

	return s.guildRepo.Delete(guildId)
}

func (s *GuildService) GetUserGuilds(userId types.ID) ([]domain.Guild, error) {
	mships, err := s.membersSvc.GetUserMemberships(userId)
	if err != nil {
		return nil, err
	}

	if len(mships) == 0 {
        return []domain.Guild{}, nil
    }

    guildIDs := make([]types.ID, len(mships))
    for i, m := range mships {
        guildIDs[i] = m.GuildId
    }

    return s.guildRepo.FindByIDs(guildIDs)
}
