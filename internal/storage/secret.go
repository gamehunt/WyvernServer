package storage

import (
	"encoding/json"
	"log/slog"
	"os"
	"wyvern/server/internal/domain"
	"log"
	"github.com/bytemare/opaque"
)

func GetSecrets(conf *opaque.Configuration) domain.Secret {
	logger := slog.Default()
	filePath := "secret.json"

	var secret domain.Secret // TODO move to config

	jsonData, err := os.ReadFile(filePath)
    if err == nil {
    	err = json.Unmarshal(jsonData, &secret)
    	if err != nil {
    	    log.Fatalf("failed to unmarshal json data: %v", err)
    	}
	} else {
    	log.Fatalf("failed to read json file: %v", err)
	}
	
	if len(secret.SecretOprfSeed) == 0 {
		logger.Info("Creating default secrets...")

		secret.SecretOprfSeed = conf.GenerateOPRFSeed()
		secret.ServerPrivateKey, secret.ServerPublicKey = conf.KeyGen()

		jsonData, err := json.Marshal(secret)
		if err != nil {
			log.Fatalf("failed to marshal json data: %v", err)
		}

		err = os.WriteFile(filePath, jsonData, 0600)
    	if err != nil {
    	    log.Fatalf("Error writing to file: %v", err)
    	}
	}

	return secret
}
