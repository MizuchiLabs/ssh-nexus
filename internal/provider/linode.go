package provider

import (
	"context"
	"net/http"
	"strconv"

	"github.com/linode/linodego"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/oauth2"
)

// LinodeProvider represents a provider implementation for Linode
type LinodeProvider struct {
	Config *models.Record
}

// NewLinodeProvider creates a new Linode provider
func NewLinodeProvider(config *models.Record) *LinodeProvider {
	return &LinodeProvider{Config: config}
}

// Sync returns a list of machines
func (p *LinodeProvider) Sync() ([]ProviderMachine, error) {
	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.Config.GetString("token")})
	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	client := linodego.NewClient(oauth2Client)
	instances, err := client.ListInstances(ctx, nil)
	if err != nil {
		return nil, err
	}

	var machines []ProviderMachine
	for _, instance := range instances {
		machine := ProviderMachine{
			ID:      strconv.Itoa(instance.ID),
			Host:    string(*instance.IPv4[0]),
			Name:    instance.Label,
			Running: instance.Status == "running",
		}
		machines = append(machines, machine)
	}

	return machines, nil
}
