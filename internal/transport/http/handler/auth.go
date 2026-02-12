package handler

import (
	"log/slog"
	"net/http"
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
	Logger         *slog.Logger
	AuthService    *service.AuthService 
	SessionService *service.SessionService
}

func NewAuthHandler(logger *slog.Logger, 
	authService *service.AuthService, 
	sessionService *service.SessionService) *AuthHandler {
	return &AuthHandler{
		Logger: logger,
		AuthService: authService,
		SessionService: sessionService,
	}
}

func (r *AuthHandler) Login(c *gin.Context)  {
	var req LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	if req.Identity != nil {
		ke2, identity, opaqueId, err := r.AuthService.StartLogin(*req.Identity, req.Payload)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			resp := LoginResponse{
				Identity:   identity,
				OpaqueId:   opaqueId,
				Payload:    ke2,
			}
			c.JSON(http.StatusOK, resp)
		}
	} else if req.OpaqueId != nil {
		sessionKey, userId, err := r.AuthService.FinishLogin(*req.OpaqueId, req.Payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		session, refreshToken, accessToken, err := r.SessionService.NewSession(*userId, sessionKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		loginResponse := SuccessLoginResponse{
			SessionId:    session.Id,
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
		}

		c.JSON(http.StatusOK, loginResponse)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
	}
}

func (r *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	if req.Identity == nil {
		payload, identity, userId, err := r.AuthService.StartRegister(req.Payload)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := RegisterResponse{
			Identity:  identity,
			UserId:   *userId,
			Payload:   payload,
		}
		
		c.JSON(http.StatusOK, resp)
	} else {
		err := r.AuthService.FinishRegister(*req.UserId, *req.Identity, req.Payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})	
	}
}

func (r *AuthHandler) Refresh(c* gin.Context) {
	var req RefreshRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	accessToken, err := r.SessionService.RefreshSession(req.SessionId, req.RefreshToken, req.Proof, req.Nonce) 

	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	resp := RefreshResponse{
		AccessToken: *accessToken,
	}

	c.JSON(http.StatusOK, resp)
}
