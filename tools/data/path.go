// Package data implements utility functions for setting fixed paths and
// various functions for creating SSH/TLS keys and certificates
package data

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Path(rel string) string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	nexusDir := filepath.Join(configDir, "ssh-nexus")
	certsDir := filepath.Join(configDir, "ssh-nexus", "certs")
	if err = os.MkdirAll(nexusDir, 0700); err != nil {
		log.Fatal(err)
	}
	if err = os.MkdirAll(certsDir, 0700); err != nil {
		log.Fatal(err)
	}

	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(nexusDir, rel)
}

// ChecksumFile calculates the sha256 checksum of a file
func ChecksumFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Checksum calculates the sha256 checksum of a string
func Checksum(data string) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(data)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
