package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"wyvern/server/internal/service"
	"wyvern/server/internal/transport/http/handler"
)

func addRoutes(mux *http.ServeMux, logger *slog.Logger, authSvc *service.AuthService) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK");
	})

	mux.HandleFunc("POST /login", handler.LoginHandler(logger, authSvc))
}

func NewRouter(
	logger *slog.Logger,
	authSvc *service.AuthService,
) http.Handler {
	r := http.NewServeMux()

	addRoutes(r, logger, authSvc)

	return r
}
