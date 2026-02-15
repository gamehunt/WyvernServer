package handler

import (
	"log/slog"
	"net/http"
	"wyvern/server/internal/pkg/errors"
	"wyvern/server/internal/pkg/util"
	"wyvern/server/internal/service"
	"wyvern/server/internal/types"

	"github.com/gin-gonic/gin"
)

type GuildCreateRequest struct {
	Name string `json:"name"`
}

type GuildHandler struct {
	logger        *slog.Logger
	guildService  *service.GuildService 
	memberService *service.MemberService
}

func NewGuildHandler(logger *slog.Logger, 
guilds *service.GuildService,
members *service.MemberService) *GuildHandler {
	return &GuildHandler{
		logger: logger,
		guildService: guilds,
		memberService: members,
	}
}

func (r *GuildHandler) GetGuild(c *gin.Context)  {
	guildId, err := types.ParseID(c.Param("guildId"))
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	guild, err := r.guildService.GetGuild(guildId)
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, guild)
}

func (r *GuildHandler) CreateGuild(c *gin.Context)  {
	var req GuildCreateRequest

    if err := c.ShouldBindJSON(&req); err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
        return
    }

	userIdParam, exists := c.Get("userId")
	if !exists {
		util.HttpError(c, http.StatusUnauthorized, errors.InvalidSession)
		return
	}

	userIdStr, ok := userIdParam.(string)
	if !ok {
		util.HttpError(c, http.StatusUnauthorized, errors.InvalidSession)
		return
	}

	userId, err := types.ParseID(userIdStr)
	if err != nil {
		util.HttpError(c, http.StatusUnauthorized, err)
        return
	}

	guild, err := r.guildService.CreateGuild(userId, req.Name)
	if err != nil {
		util.HttpError(c, http.StatusInternalServerError, err)
        return
	}

	c.JSON(http.StatusOK, guild)
}

func (r *GuildHandler) DeleteGuild(c *gin.Context)  {
	guildId, err := types.ParseID(c.Param("guildId"))
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	requesterId, err := types.ParseID(c.Param("userId"))
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	err = r.guildService.DeleteGuild(requesterId, guildId)
	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusNoContent)
}
