package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"wyvern/server/internal/service"
	"wyvern/server/internal/types"
)

type LoginRequest struct {
	Identity  *string `json:"identity,omitempty"`
	OpaqueId  *types.ID `json:"opaque_id,omitempty"`
	Payload   []byte `json:"payload"`
}

type LoginResponse struct {
	Identity  []byte   `json:"identity,omitempty"`
	OpaqueId *types.ID `json:"opaque_id,omitempty"`
	Payload   []byte
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

func LoginHandler(_ *slog.Logger, authService *service.AuthService, sessionService *service.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)

		var req LoginRequest

		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Identity != nil {
			ke2, identity, opaqueId, err := authService.StartLogin(*req.Identity, req.Payload)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				resp := LoginResponse{
					Identity:   identity,
					OpaqueId:   opaqueId,
					Payload:    ke2,
				}

				err = json.NewEncoder(w).Encode(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		} else if req.OpaqueId != nil {
			_, userId, err := authService.FinishLogin(*req.OpaqueId, req.Payload)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			session, err := sessionService.NewSession(*userId)

			loginResponse := SuccessLoginResponse{
				SessionId:    session.Id,
				AccessToken:  "",
				RefreshToken: "",
			}

			err = json.NewEncoder(w).Encode(loginResponse)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}
	}
}

func RegisterHandler(log *slog.Logger, authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)

		var req RegisterRequest

		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Identity == nil {
			payload, identity, userId, err := authService.StartRegister(req.Payload)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			req := RegisterResponse{
				Identity:  identity,
				UserId:   *userId,
				Payload:   payload,
			}
			
			err = json.NewEncoder(w).Encode(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := authService.FinishRegister(*req.UserId, *req.Identity, req.Payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		}
	}
}
