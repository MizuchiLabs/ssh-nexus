package provider

import (
	"github.com/pocketbase/pocketbase/models"
)

type AzureProvider struct {
	Config *models.Record
}

func NewAzureProvider(config *models.Record) *AzureProvider {
	return &AzureProvider{Config: config}
}

func (p *AzureProvider) Sync() ([]ProviderMachine, error) {
	// TODO: Implement

	var machines []ProviderMachine

	return machines, nil
}
