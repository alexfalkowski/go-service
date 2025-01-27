package token

import (
	"github.com/agentstation/uuidkey"
	"github.com/alexfalkowski/go-service/os"
	"github.com/google/uuid"
)

// GenerateKey for token.
func GenerateKey() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	key, err := uuidkey.NewAPIKey(os.ExecutableName(), id.String())
	if err != nil {
		return "", err
	}

	return key.String(), nil
}

// VerifyKey for token.
func VerifyKey(key string) error {
	_, err := uuidkey.ParseAPIKey(key)

	return err
}
