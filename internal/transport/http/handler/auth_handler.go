package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"wyvern/server/internal/service"
)

type loginRequest struct {
	Identity *[]byte `json:"identity,omitempty"`
	Payload   []byte `json:"payload"`
}

type loginChallengeResponse struct {
	Payload []byte
}

type loginSuccessResponse struct {
	SessionId []byte `json:"token"`
}

func LoginHandler(_ *slog.Logger, authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)

		var req loginRequest

		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if req.Identity != nil {
			ke2, err := authService.StartLogin(*req.Identity, req.Payload)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				resp := loginChallengeResponse{
					Payload: ke2,
				}

				err = json.NewEncoder(w).Encode(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		} else {
			sessionId, err := authService.FinishLogin(req.Payload)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			loginResponse := loginSuccessResponse{
				SessionId: *sessionId,
			}

			err = json.NewEncoder(w).Encode(loginResponse)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	}
}
