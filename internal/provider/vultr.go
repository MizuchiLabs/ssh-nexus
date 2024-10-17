package provider

import (
	"context"

	"github.com/pocketbase/pocketbase/models"
	"github.com/vultr/govultr/v3"
	"golang.org/x/oauth2"
)

type VultrProvider struct {
	Config *models.Record
}

func NewVultrProvider(config *models.Record) *VultrProvider {
	return &VultrProvider{Config: config}
}

func (p *VultrProvider) Sync() ([]ProviderMachine, error) {
	config := &oauth2.Config{}
	ctx := context.Background()
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: p.Config.GetString("token")})
	client := govultr.NewClient(oauth2.NewClient(ctx, ts))

	var machines []ProviderMachine
	listOptions := &govultr.ListOptions{PerPage: 100}
	for {
		instances, meta, _, err := client.Instance.List(ctx, listOptions)
		if err != nil {
			return nil, err
		}
		for _, v := range instances {
			machine := ProviderMachine{
				ID:      v.ID,
				Host:    v.MainIP,
				Name:    v.Label,
				Running: v.Status == "active",
			}
			machines = append(machines, machine)
		}

		if meta.Links.Next == "" {
			break
		} else {
			listOptions.Cursor = meta.Links.Next
			continue
		}
	}

	return machines, nil
}
