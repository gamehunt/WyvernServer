package service

import (
	"wyvern/server/internal/domain"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"
)

type PermissionService struct {
	permOverrideRepo repository.PermissionOverrideRepository	
	membersRepo repository.MemberRepository
	rolesRepo   repository.RoleRepository
}

func NewPermissionService(
	permsOverrideRepo repository.PermissionOverrideRepository,
	memberRepo repository.MemberRepository,
	roleRepo   repository.RoleRepository,
) *PermissionService {
	return &PermissionService {
		permOverrideRepo: permsOverrideRepo,
		membersRepo: memberRepo,
		rolesRepo: roleRepo,
	}	
}

func (s *PermissionService) ComputeGuildPermissions(
    userId, guildId types.ID,
) (domain.Permission, error) {

    member, err := s.membersRepo.GetForUser(guildId, userId)
    if err != nil {
        return 0, err
    }

    roles, err := s.rolesRepo.FindByIDs(member.Roles)
    if err != nil {
        return 0, err
    }

    var perms domain.Permission = 0

    for _, role := range roles {
        perms = perms.Add(role.Permissions)
    }

    return perms, nil
}

func (s *PermissionService) ComputeChannelPermissions(
    userId, guildId, channelId types.ID,
) (domain.Permission, error) {

    perms, err := s.ComputeGuildPermissions(userId, guildId)
    if err != nil {
        return 0, err
    }

    overwrites, err := s.permOverrideRepo.FindByChannel(channelId)
    if err != nil {
        return 0, err
    }

    member, err := s.membersRepo.GetForUser(userId, guildId)
    if err != nil {
        return 0, err
    }

	// everyone
    for _, o := range overwrites {
        if o.Type == domain.RoleOverride && o.TargetId == guildId {
            perms = perms.Remove(o.Deny)
            perms = perms.Add(o.Allow)
        }
    }

    for _, roleId := range member.Roles {
        for _, o := range overwrites {
            if o.Type == domain.RoleOverride && o.TargetId == roleId {
                perms = perms.Remove(o.Deny)
                perms = perms.Add(o.Allow)
            }
        }
    }

    for _, o := range overwrites {
        if o.Type == domain.MemberOverride && o.TargetId == userId {
            perms = perms.Remove(o.Deny)
            perms = perms.Add(o.Allow)
        }
    }

    return perms, nil
}

func (s *PermissionService) CanManageGuild(userId, guildId types.ID) (bool, error) {
	perms, err := s.ComputeGuildPermissions(userId, guildId)
	if err != nil {
		return false, err
	}
	return perms.Has(domain.ManageGuild) || perms.Has(domain.Administrator), nil
}  
