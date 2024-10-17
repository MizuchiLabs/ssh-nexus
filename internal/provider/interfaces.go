// Package provider contains all the external provider implementations
package provider

import (
	"fmt"

	"github.com/pocketbase/pocketbase/models"
)

type Provider interface {
	Sync() ([]ProviderMachine, error)
}

type ProviderMachine struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Host    string `json:"host,omitempty"`
	Running bool   `json:"running,omitempty"`
}

// NewProvider creates a new provider
func NewProvider(config *models.Record) (Provider, error) {
	username := config.GetString("username")
	password := config.GetString("password")
	token := config.GetString("token")

	switch config.GetString("type") {
	case "aws":
		if username == "" || password == "" || token == "" {
			return nil, fmt.Errorf("please provide an id, secret and token")
		}
		return NewAWSProvider(config), nil
	case "linode":
		if token == "" {
			return nil, fmt.Errorf("please provide a token")
		}
		return NewLinodeProvider(config), nil
	case "hetzner":
		if token == "" {
			return nil, fmt.Errorf("please provide a token")
		}
		return NewHetznerProvider(config), nil
	case "vultr":
		if token == "" {
			return nil, fmt.Errorf("please provide a token")
		}
		return NewVultrProvider(config), nil
	case "proxmox":
		if username == "" {
			return nil, fmt.Errorf("please provide a username")
		}
		if token == "" && password == "" {
			return nil, fmt.Errorf("please provide either a token or password")
		}
		return NewProxmoxProvider(config), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.GetString("type"))
	}
}
