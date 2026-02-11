package storage

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"wyvern/server/internal/domain"

	"github.com/bytemare/opaque"
)

func GetSecrets(filePath string, conf *opaque.Configuration) domain.Secret {
	logger := slog.Default()

	var secret domain.Secret 

	_, err := os.Stat(filePath)

	if err == nil {
		jsonData, err := os.ReadFile(filePath)
    	err = json.Unmarshal(jsonData, &secret)
    	if err != nil {
    	    log.Fatalf("failed to unmarshal json data: %v", err)
    	}
	} else if errors.Is(err, os.ErrNotExist) {
		logger.Info("Creating default secrets...")

		secret.SecretOprfSeed = conf.GenerateOPRFSeed()
		secret.ServerPrivateKey, secret.ServerPublicKey = conf.KeyGen()

		jsonData, err := json.MarshalIndent(secret, "", "\t")
		if err != nil {
			log.Fatalf("failed to marshal json data: %v", err)
		}

		err = os.WriteFile(filePath, jsonData, 0600)
    	if err != nil {
    	    log.Fatalf("Error writing to file: %v", err)
    	}
	} else {
    	log.Fatalf("failed to read json file: %v", err)
	}
	
	return secret
}
