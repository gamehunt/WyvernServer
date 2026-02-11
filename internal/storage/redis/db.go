package redis

import (
	"github.com/valkey-io/valkey-go"
)

func New(uri string) (*valkey.Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{uri}})

	if err != nil {
		return nil, err
	}

	return &client, nil
}
