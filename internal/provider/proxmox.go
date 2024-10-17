package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/luthermonson/go-proxmox"
	"github.com/pocketbase/pocketbase/models"
)

type ContainerInterface struct {
	Hwaddr string `json:"hwaddr,omitempty"`
	Name   string `json:"name,omitempty"`
	Inet   string `json:"inet,omitempty"`
	Inet6  string `json:"inet6,omitempty"`
}

type ProxmoxProvider struct {
	Config *models.Record
}

func NewProxmoxProvider(config *models.Record) *ProxmoxProvider {
	return &ProxmoxProvider{Config: config}
}

func (p *ProxmoxProvider) Sync() ([]ProviderMachine, error) {
	ctx := context.Background()

	opts := []proxmox.Option{
		proxmox.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}),
	}

	if p.Config.GetString("token") != "" {
		opts = append(opts, proxmox.WithAPIToken(
			p.Config.GetString("username"),
			p.Config.GetString("token"),
		))
	} else {
		opts = append(opts, proxmox.WithCredentials(&proxmox.Credentials{
			Username: p.Config.GetString("username"),
			Password: p.Config.GetString("password"),
		}))
	}

	u, err := url.Parse(p.Config.GetString("url"))
	if err != nil {
		return nil, err
	}

	url := p.Config.GetString("url")
	if u.Scheme == "" {
		url = "https://" + url
	}

	client := proxmox.NewClient(fmt.Sprintf("%s:8006/api2/json", url), opts...)
	cluster, err := client.Cluster(ctx)
	if err != nil {
		return nil, err
	}

	resources, err := cluster.Resources(ctx)
	if err != nil {
		return nil, err
	}

	// We're only interested in VMs & Containers, sorted by node name
	sorted := make(map[string]proxmox.ClusterResources)
	for _, resource := range resources {
		if slices.Contains([]string{"qemu", "lxc"}, resource.Type) && resource.Template == 0 {
			sorted[resource.Node] = append(sorted[resource.Node], resource)
		}
	}

	// Since we need info from every node in a cluster we need to create multiple clients
	// Rather chaotic way of fetching machines, but here the info to get the ip
	// from a vm is fetched by the qemu agent and for lxc containers it's from the config
	var machines []ProviderMachine
	for node, resource := range sorted {
		nodeClient, err := getNodeClient(ctx, client, opts, node)
		if err != nil {
			slog.Error("Failed to get node client", "Error", err)
			continue
		}

		nodeInfo, err := nodeClient.Node(ctx, node)
		if err != nil {
			slog.Error("Failed to get node info", "Error", err)
			continue
		}

		vms, _ := nodeInfo.VirtualMachines(ctx)
		for _, vm := range vms {
			if proxmox.IsTemplate(vm.Template) {
				continue
			}
			ip, err := getVMHost(ctx, nodeClient, vm, node)
			if err != nil {
				slog.Error("Failed to get vm host", "Error", err)
				continue
			}

			machine := ProviderMachine{
				ID:      strconv.Itoa(int(vm.VMID)),
				Name:    vm.Name,
				Host:    ip,
				Running: ip != "" && vm.Status == "running",
			}
			machines = append(machines, machine)
		}

		for _, r := range resource {
			if r.Type != "lxc" {
				continue
			}
			ip, err := getContainerHost(ctx, nodeClient, node, int(r.VMID))
			if err != nil {
				slog.Error("Failed to get container host", "Error", err)
				continue
			}

			machine := ProviderMachine{
				ID:      strconv.Itoa(int(r.VMID)),
				Name:    r.Name,
				Host:    ip,
				Running: ip != "" && r.Status == "running",
			}
			machines = append(machines, machine)
		}
	}

	return machines, nil
}

func getNodeClient(
	ctx context.Context,
	client *proxmox.Client,
	opts []proxmox.Option,
	node string,
) (*proxmox.Client, error) {
	var network []struct {
		Address string `json:"address,omitempty"`
		Gateway string `json:"gateway,omitempty"`
		Iface   string `json:"iface,omitempty"`
		Netmask string `json:"netmask,omitempty"`
		Type    string `json:"type,omitempty"`
	}

	err := client.Get(ctx, fmt.Sprintf("/nodes/%s/network", node), &network)
	if err != nil {
		return nil, err
	}

	for _, net := range network {
		if net.Address != "" {
			return proxmox.NewClient(
				fmt.Sprintf("https://%s:8006/api2/json", net.Address),
				opts...), nil
		}
	}

	return nil, nil
}

func getVMHost(
	ctx context.Context,
	client *proxmox.Client,
	vm *proxmox.VirtualMachine,
	node string,
) (string, error) {
	ifaces, err := vm.AgentGetNetworkIFaces(ctx)
	if err != nil {
		return "", err
	}
	if len(ifaces) == 0 {
		return "", nil
	}

	var data map[string]any
	err = client.Get(ctx, fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vm.VMID), &data)
	if err != nil {
		return "", err
	}

	// We need to compare the mac address of the virtio interface to the one
	// provided by the qemu agent
	var macAddr string
	for k, v := range data {
		if strings.HasPrefix(k, "net") {
			parts := strings.Split(v.(string), ",")
			macAddr = strings.TrimPrefix(parts[0], "virtio=")
			break
		}
	}

	for _, iface := range ifaces {
		if iface.HardwareAddress == strings.ToLower(macAddr) {
			for _, ip := range iface.IPAddresses {
				if ip.IPAddressType == "ipv4" {
					return strings.Split(ip.IPAddress, "/")[0], nil
				}
			}
		}
	}

	return "", nil
}

// Attempt to get first ip address in the list
// TODO: Maybe handle multiple ip addresses
func getContainerHost(
	ctx context.Context,
	client *proxmox.Client,
	node string,
	vmid int,
) (string, error) {
	var data []ContainerInterface
	err := client.Get(ctx, fmt.Sprintf("/nodes/%s/lxc/%d/interfaces", node, vmid), &data)
	if err != nil {
		return "", err
	}

	for _, iface := range data {
		if iface.Inet == "" || strings.Contains(iface.Inet, "127.0.0.1") {
			continue
		}
		return strings.Split(iface.Inet, "/")[0], nil
	}
	return "", nil
}
