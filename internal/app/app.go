package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	nethttp "net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"wyvern/server/internal/config"
	"wyvern/server/internal/pkg/logger"
	"wyvern/server/internal/service"
	mongoimpl "wyvern/server/internal/storage/mongo"
	transporthttp "wyvern/server/internal/transport/http"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


type App struct {
	httpServer *nethttp.Server
	// wsManager  *ws.ConnectionManager
	mongoClient *mongo.Client
}

func New(cfg config.Config) (*App, error) {
	log := logger.New(cfg.Logger)
	logger.SetAsDefault(log)

	mongoClient, err := mongoimpl.New(cfg.Mongo.Uri)
	if err != nil {
		log.Error("Failed to connect to MongoDB", "error", err)
		return nil, err
	}

	userRepo := mongoimpl.NewUserRepository(mongoClient, cfg.Mongo.Database)

	authSvc := service.NewAuthService(userRepo, cfg.Auth)

	handler := transporthttp.NewRouter(log, authSvc)

	httpServer := &nethttp.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
	}

	return &App{
		httpServer: httpServer,
		mongoClient: mongoClient,
	}, nil
}

func (a *App) Run() error {
	log := slog.Default()

	go func() {
		log.Info("Listening at", "addr", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
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

	if err := mongoimpl.CloseGracefully(a.mongoClient); err != nil {
		log.Error("Failed to shutdown MongoDB", "error", err)
		return fmt.Errorf("Failed to shutdown MongoDB: %w", err)
	}

	log.Info("Bye.")
	return nil
}
