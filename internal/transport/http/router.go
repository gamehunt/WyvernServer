package http

import (
	"fmt"
	"log/slog"
	"net/http"
)

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test");
	})
}

func NewRouter(
	logger *slog.Logger,
) http.Handler {
	r := http.NewServeMux()

	addRoutes(r)

	return r
}
