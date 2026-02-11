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
	"wyvern/server/internal/storage/redis"
	transporthttp "wyvern/server/internal/transport/http"

	"github.com/valkey-io/valkey-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


type App struct {
	httpServer *nethttp.Server
	// wsManager  *ws.ConnectionManager
	mongoClient *mongo.Client
	redisClient valkey.Client
}

func New(cfg config.Config) (*App, error) {
	log := logger.New(cfg.Logger)
	logger.SetAsDefault(log)

	mongoClient, err := mongoimpl.New(cfg.Mongo)
	if err != nil {
		log.Error("Failed to connect to MongoDB", "error", err)
		return nil, err
	}

	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		log.Error("Failed to connect to Redis", "error", err)
		return nil, err
	}

	userRepo := mongoimpl.NewUserRepository(mongoClient, cfg.Mongo.Database)

	opaqueSessionStorate := redis.NewOpaqueSessionStorage(redisClient)

	authSvc := service.NewAuthService(userRepo, opaqueSessionStorate, cfg.Auth)

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
		redisClient: redisClient,
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

	a.redisClient.Close()

	log.Info("Bye.")
	return nil
}
