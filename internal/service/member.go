package service

import (
	"time"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/pkg/errors"
	"wyvern/server/internal/repository"
	"wyvern/server/internal/types"
)

type MemberService struct {
	memberRepo repository.MemberRepository
	roleRepo   repository.RoleRepository
}

func NewMemberService(
	memberRepo repository.MemberRepository,
	roleRepo   repository.RoleRepository,
) *MemberService {
	return &MemberService {
		memberRepo: memberRepo,
		roleRepo: roleRepo,
	}	
}

func (s *MemberService) JoinGuild(userId, guildId types.ID) (*domain.Member, error) {
	member := &domain.Member{
		UserId: userId,
		GuildId: guildId,
		JoinedAt: time.Now(),
	}

	err := s.memberRepo.Create(member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (s *MemberService) LeaveGuild(userId, guildId types.ID) error {
	return s.memberRepo.Delete(userId, guildId)
}

func (s *MemberService) GetMember(guildId, userId types.ID) (*domain.Member, error) {
	return s.memberRepo.GetForUser(userId, guildId)
}

func (s *MemberService) ListMembers(guildId types.ID) ([]domain.Member, error) {
	return s.memberRepo.FindByGuild(guildId)
}

func (s *MemberService) GetMemberRoles(guildId, userId types.ID) ([]domain.Role, error){
	member, err := s.GetMember(guildId, userId)
	if err != nil {
		return nil, err
	}

	if member == nil {
		return nil, errors.InvalidMember
	}

	return s.roleRepo.FindByIDs(member.Roles)
}

func (s *MemberService) GetUserMemberships(userId types.ID) ([]domain.Member, error){
	return s.memberRepo.FindByUser(userId)
}
