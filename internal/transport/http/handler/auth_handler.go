package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"wyvern/server/internal/service"
)

type LoginRequest struct {
	Identity  *string `json:"identity,omitempty"`
	SessionId *string `json:"session_id,omitempty"`
	Payload   []byte `json:"payload"`
}

type LoginResponse struct {
	Identity  []byte `json:"identity,omitempty"`
	SessionId *string `json:"session_id,omitempty"`
	Payload   []byte
}

type RegisterRequest struct {
	Identity *string `json:"identity,omitempty"`
	CredId   *string `json:"cred_id,omitempty"`
	Payload   []byte `json:"payload"`
}

type RegisterResponse struct {
	Identity []byte `json:"identity"`
	CredId   string `json:"cred_id"`
	Payload  []byte `json:"payload"`
}

func LoginHandler(_ *slog.Logger, authService *service.AuthService) http.HandlerFunc {
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
			ke2, identity, sessionId, err := authService.StartLogin(*req.Identity, req.Payload)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				resp := LoginResponse{
					Identity:   identity,
					SessionId: &sessionId,
					Payload:    ke2,
				}

				err = json.NewEncoder(w).Encode(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		} else if req.SessionId != nil {
			_, err := authService.FinishLogin(*req.SessionId, req.Payload)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// loginResponse := LoginResponse{
			// 	Payload: sessionId,
			// }
			//
			// err = json.NewEncoder(w).Encode(loginResponse)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// }

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		} else {
			http.Error(w, "Invalid reques", http.StatusBadRequest)
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
			payload, identity, credId, err := authService.StartRegister(req.Payload)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			req := RegisterResponse{
				Identity: identity,
				CredId:   credId,
				Payload:  payload,
			}
			
			err = json.NewEncoder(w).Encode(req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			err := authService.FinishRegister(*req.CredId, *req.Identity, req.Payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		}
	}
}
