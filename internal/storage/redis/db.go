package redis

import (
	"wyvern/server/internal/config"

	"github.com/valkey-io/valkey-go"
)

func New(cfg config.RedisConfig) (valkey.Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.Uri}, SelectDB: cfg.DB})

	if err != nil {
		return nil, err
	}

	return client, nil
}
