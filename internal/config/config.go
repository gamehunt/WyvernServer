package config

import (
	"os"
	"strconv"
	"wyvern/server/internal/pkg/logger"
)

type Config struct {
	Server ServerConfig
	Mongo  MongoConfig
	Redis  RedisConfig
	Logger logger.LoggerConfig
}

type ServerConfig struct {
	Host string
	Port int
	Timeout int
}

type RedisConfig struct {
	Host     string
	Port     int
	DB       int
}

type MongoConfig struct {
	Uri      string
	Database string
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func Load() (*Config) {
	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", ""),
			Port: getEnvInt("SERVER_PORT", 5000),
			Timeout: getEnvInt("SERVER_TIMEOUT", 5000),
		},
		Mongo: MongoConfig{
			Uri:      getEnv("MONGO_URI", "mongodb://wyvern:wyvern@mongo:27017"),
			Database: getEnv("MONGO_DATABASE", "wyvern"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "valkey"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Logger: logger.LoggerConfig{
			Level:    logger.LevelFromString(getEnv("LOG_LEVEL", "info")),
		},
	}

	return cfg
}
