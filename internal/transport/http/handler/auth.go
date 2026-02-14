package handler

import (
	"log/slog"
	"net/http"
	"wyvern/server/internal/pkg/util"
	"wyvern/server/internal/service"
	"wyvern/server/internal/types"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Identity  *string   `json:"identity,omitempty"`
	OpaqueId  *types.ID `json:"opaque_id,omitempty"`
	Payload   []byte    `json:"payload"`
}

type LoginResponse struct {
	Identity  []byte   `json:"identity,omitempty"`
	OpaqueId *types.ID `json:"opaque_id,omitempty"`
	Payload   []byte   `json:"payload"`
}

type SuccessLoginResponse struct {
	SessionId    types.ID `json:"session_id"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

type RegisterRequest struct {
	UserId   *types.ID `json:"user_id,omitempty"`
	Identity *string   `json:"identity,omitempty"`
	Payload   []byte   `json:"payload"`
}

type RegisterResponse struct {
	UserId   types.ID `json:"user_id"`
	Identity []byte   `json:"identity"`
	Payload  []byte   `json:"payload"`
}

type RefreshRequest struct {
	SessionId     types.ID `json:"session_id"`
	RefreshToken  string   `json:"refresh_token"`
	Proof         []byte   `json:"proof"`
	Nonce         []byte   `json:"nonce"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type AuthHandler struct {
	logger         *slog.Logger
	authService    *service.AuthService 
	sessionService *service.SessionService
}

func NewAuthHandler(logger *slog.Logger, 
	authService *service.AuthService, 
	sessionService *service.SessionService) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		authService: authService,
		sessionService: sessionService,
	}
}

// Login         godoc
// @Summary      Logs in
// @Description  Logs in
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /auth/login [post]
func (r *AuthHandler) Login(c *gin.Context)  {
	var req LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
        return
    }

	if req.Identity != nil {
		ke2, identity, opaqueId, err := r.authService.StartLogin(*req.Identity, req.Payload)

		if err != nil {
			util.HttpError(c, http.StatusBadRequest, err)
		} else {
			resp := LoginResponse{
				Identity:   identity,
				OpaqueId:   opaqueId,
				Payload:    ke2,
			}
			c.JSON(http.StatusOK, resp)
		}
	} else if req.OpaqueId != nil {
		sessionKey, userId, err := r.authService.FinishLogin(*req.OpaqueId, req.Payload)
		if err != nil {
			util.HttpError(c, http.StatusBadRequest, err)
			return
		}

		session, refreshToken, accessToken, err := r.sessionService.NewSession(*userId, sessionKey)
		if err != nil {
			util.HttpError(c, http.StatusInternalServerError, err)
			return
		}

		loginResponse := SuccessLoginResponse{
			SessionId:    session.Id,
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
		}

		c.JSON(http.StatusOK, loginResponse)
	} else {
		util.HttpErrorString(c, http.StatusBadRequest, "Invalid request")
	}
}

// Register      godoc
// @Summary      Registers
// @Description  Registers
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /auth/register [post]
func (r *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
        return
    }

	if req.Identity == nil {
		payload, identity, userId, err := r.authService.StartRegister(req.Payload)

		if err != nil {
			util.HttpError(c, http.StatusBadRequest, err)
			return
		}

		resp := RegisterResponse{
			Identity:  identity,
			UserId:   *userId,
			Payload:   payload,
		}
		
		c.JSON(http.StatusOK, resp)
	} else {
		err := r.authService.FinishRegister(*req.UserId, *req.Identity, req.Payload)
		if err != nil {
			util.HttpError(c, http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{})	
	}
}

// Refresh       godoc
// @Summary      Refreshes
// @Description  Refreshes
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /auth/refresh [post]
func (r *AuthHandler) Refresh(c* gin.Context) {
	var req RefreshRequest
    if err := c.ShouldBindJSON(&req); err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
        return
    }

	accessToken, err := r.sessionService.RefreshSession(req.SessionId, req.RefreshToken, req.Proof, req.Nonce) 

	if err != nil {
		util.HttpError(c, http.StatusBadRequest, err)
		return 
	}

	resp := RefreshResponse{
		AccessToken: *accessToken,
	}

	c.JSON(http.StatusOK, resp)
}
