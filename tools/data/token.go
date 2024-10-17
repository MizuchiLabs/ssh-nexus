package data

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

func GenerateToken(rotate bool) error {
	if rotate {
		if err := os.Remove(Token); err != nil {
			return err
		}
	}

	if _, err := os.Stat(Token); err == nil {
		return nil
	}

	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return err
	}

	token := hex.EncodeToString(bytes)
	err = os.WriteFile(Token, []byte(token), 0600)
	if err != nil {
		return err
	}
	return nil
}

func GetToken() (string, error) {
	token, err := os.ReadFile(Token)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
