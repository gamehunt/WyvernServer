package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"wyvern/server/internal/pkg/logger"
)

type Config struct {
	Server ServerConfig        `json:"server"`
	Auth   AuthConfig          `json:"auth"`
	Mongo  MongoConfig         `json:"mongo"`
	Redis  RedisConfig         `json:"redis"`
	Logger logger.LoggerConfig `json:"logger"`
}

type ServerConfig struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Timeout int    `json:"timeout"`
}

type AuthConfig struct {
	ServerIdentity string `json:"server_identity"`
	SecretPath     string `json:"secret_path"`
}

type RedisConfig struct {
	Uri string `json:"uri"`
	DB  int    `json:"db"`
}

type MongoConfig struct {
	Uri      string `json:"uri"`
	Database string `json:"database"`
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

func Load(path string) (Config) {
	var cfg Config
	
	jsonData, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	err = json.Unmarshal(jsonData, &cfg)

	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return cfg
}
