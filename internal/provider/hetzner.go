package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/pocketbase/pocketbase/models"
)

type HetznerProvider struct {
	Config *models.Record
}

func NewHetznerProvider(config *models.Record) *HetznerProvider {
	return &HetznerProvider{Config: config}
}

func (p *HetznerProvider) Sync() ([]ProviderMachine, error) {
	client := hcloud.NewClient(
		hcloud.WithToken(p.Config.GetString("token")),
		hcloud.WithBackoffFunc(func(retries int) time.Duration {
			return time.Duration(retries*5) * time.Second
		}),
	)
	servers, err := client.Server.All(context.Background())
	if err != nil {
		return nil, err
	}

	var machines []ProviderMachine
	for _, server := range servers {
		machine := ProviderMachine{
			ID:      strconv.FormatInt(server.ID, 10),
			Host:    server.PublicNet.IPv4.IP.String(),
			Name:    server.Name,
			Running: server.Status == "running",
		}
		machines = append(machines, machine)
	}

	return machines, nil
}
