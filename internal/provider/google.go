package provider

import (
	"github.com/pocketbase/pocketbase/models"
)

type GoogleCloudProvider struct {
	Config *models.Record
}

func NewGoogleCloudProvider(config *models.Record) *GoogleCloudProvider {
	return &GoogleCloudProvider{Config: config}
}

func (p *GoogleCloudProvider) Sync() ([]ProviderMachine, error) {
	// TODO: Implement

	// var machines []ProviderMachine

	return nil, nil
}
