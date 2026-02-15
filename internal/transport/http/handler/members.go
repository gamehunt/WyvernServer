package handler

import (
	"log/slog"
	"net/http"
	"time"
	"wyvern/server/internal/domain"
	"wyvern/server/internal/pkg/errors"
	"wyvern/server/internal/pkg/util"
	"wyvern/server/internal/service"
	"wyvern/server/internal/types"

	"github.com/gin-gonic/gin"
)

type MemberObject struct {
	User     domain.User `json:"user"`
	Nickname *string     `json:"nickname,omitempty"`
	Roles   []types.ID   `json:"roles"`
	JoinedAt time.Time   `json:"joined_at"`
}

type MemberHandler struct {
	logger        *slog.Logger
	userService   *service.UserService
	memberService *service.MemberService
}

func NewMemberHandler(logger *slog.Logger, 
userService   *service.UserService,
members *service.MemberService) *MemberHandler {
	return &MemberHandler{
		logger: logger,
		userService: userService,
		memberService: members,
	}
}

func (r *MemberHandler) makeMemberObject(member *domain.Member) (*MemberObject, error) {
	user, err := r.userService.GetUser(member.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.InvalidMember
	}

	memberObj := &MemberObject{
		User: *user,
		Nickname: member.Nickname,
		Roles: member.Roles,
		JoinedAt: member.JoinedAt,
	}

	return memberObj, nil
}

func (r *MemberHandler) ListMembers(c *gin.Context) {
	guildId, err := types.ParseID(c.Param("guildId"))
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	members, err := r.memberService.ListMembers(guildId)
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	var memberObjs []MemberObject
    for _, v := range members {
		memberObj, err := r.makeMemberObject(&v)
		if err != nil {
			util.HttpError(c, http.StatusBadRequest, err)
			return
		}
		memberObjs = append(memberObjs, *memberObj)
    }

	c.JSON(http.StatusOK, memberObjs)
}

func (r *MemberHandler) GetMember(c *gin.Context) {
	guildId, err := types.ParseID(c.Param("guildId"))
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	userId, err := types.ParseID(c.Param("guildId"))
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	member, err := r.memberService.GetMember(guildId, userId)
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	if member == nil {
		util.HttpError(c, http.StatusBadRequest, errors.InvalidMember)
		return
	}

	memberObj, err := r.makeMemberObject(member)
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}


	c.JSON(http.StatusOK, memberObj)
}
