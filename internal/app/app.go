package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"wyvern/server/internal/config"
	"wyvern/server/internal/pkg/logger"
	"wyvern/server/internal/service"
	"wyvern/server/internal/storage/mongo"
	my_http "wyvern/server/internal/transport/http"
)


type App struct {
	httpServer *http.Server
	// wsManager  *ws.ConnectionManager
}

func New(cfg *config.Config) (*App, error) {
	log := logger.New(cfg.Logger)
	logger.SetAsDefault(log)

	mongoClient, err := mongo.New(cfg.Mongo.Uri)
	if err != nil {
		log.Error("Ошибка подключения к MongoDB", "error", err)
		return nil, err
	}

	database := mongoClient.Database(cfg.Mongo.Database)

	userRepo := mongo.NewUserRepository(database)

	service.NewAuthService(userRepo)

	handler := my_http.NewRouter(log)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
	}

	return &App{
		httpServer: httpServer,
	}, nil
}

func (a *App) Run() error {
	log := slog.Default()

	go func() {
		log.Info("Listening at", "addr", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("SIGINT")

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	log := slog.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Error("Graceful shutdown error", "error", err)
		return fmt.Errorf("Failed to shutdown server: %w", err)
	}

	log.Info("Bye.")
	return nil
}
