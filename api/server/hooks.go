package server

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	agentv1 "github.com/MizuchiLabs/ssh-nexus/api/proto/gen/agent/v1"
	"github.com/MizuchiLabs/ssh-nexus/tools/data"
	"github.com/pocketbase/pocketbase/core"
)

func BoolPointer(b bool) *bool {
	return &b
}

func (s *AgentServer) initClient(client Client) {
	response := &agentv1.StreamResponse{}
	sshConfig, err := s.PB.Dao().FindFirstRecordByData("settings", "key", "ssh_config")
	if err != nil {
		slog.Error("failed to get ssh config", "err", err)
	}

	userCa, err := data.GetPublicUserKey()
	if err != nil {
		slog.Error("failed to get user ca", "err", err)
	}

	principals, err := getPrincipals(s.PB, client.Machine)
	if err != nil {
		slog.Error("failed to get principals", "err", err)
	}

	response.SshConfig = []byte(sshConfig.GetString("value"))
	response.UserCertificatePublicKey = userCa
	response.Principals = principals

	if err := client.Stream.Send(response); err != nil {
		slog.Error("initializing agent error", "err", err)
	}
}

func (s *AgentServer) monitorHook(client Client, req *agentv1.StreamRequest) {
	response := &agentv1.StreamResponse{}

	if req.GetPublicHostKey() != "" {
		cert, err := data.SignHostCertificate(
			req.GetPublicHostKey(),
			client.Machine.GetString("name"),
			30*24*time.Hour,
		)
		if err != nil {
			slog.Error("failed to sign host cert", "err", err)
		}
		response.HostCertificatePublicKey = cert
	}
	if err := client.Stream.Send(response); err != nil {
		slog.Error("initializing agent error", "err", err)
	}
}

func (s *AgentServer) pbHook() {
	s.PB.OnRecordAfterUpdateRequest("settings").
		Add(func(e *core.RecordUpdateEvent) error {
			if e.Record.GetString("key") == "ssh_config" {
				for id := range s.Clients {
					client := s.Clients[id]
					reply := &agentv1.StreamResponse{
						SshConfig: []byte(e.Record.GetString("value")),
					}
					if err := client.Stream.Send(reply); err != nil {
						slog.Error("updating agent error", "err", err)
					}
				}
			}
			return nil
		})

	s.PB.OnRecordAfterUpdateRequest("machines").
		Add(func(e *core.RecordUpdateEvent) error {
			client := s.Clients[e.Record.Id]
			principals, err := getPrincipals(s.PB, e.Record)
			if err != nil {
				slog.Error("failed to get principals", "err", err)
			}

			reply := &agentv1.StreamResponse{Principals: principals}
			if err := client.Stream.Send(reply); err != nil {
				slog.Error("updating agent error", "err", err)
			}
			return nil
		})

	s.PB.OnRecordAfterUpdateRequest("users").Add(func(e *core.RecordUpdateEvent) error {
		// Update machines if groups changed
		groups := e.Record.GetStringSlice("groups")
		oldGroups := e.Record.OriginalCopy().GetStringSlice("groups")

		// or principal id changes
		principal := e.Record.GetString("principal")
		oldPrincipal := e.Record.OriginalCopy().GetString("principal")

		if !slices.Equal(groups, oldGroups) || (principal != oldPrincipal) {
			machinesBefore, err := getUserMachines(s.PB, e.Record.OriginalCopy())
			if err != nil {
				return err
			}

			machinesAfter, err := getUserMachines(s.PB, e.Record)
			if err != nil {
				return err
			}

			machines := append(machinesBefore, machinesAfter...)
			for _, machine := range machines {
				principals, err := getPrincipals(s.PB, machine)
				if err != nil {
					return fmt.Errorf("failed to get machine users: %v", err)
				}

				client := s.Clients[machine.Id]
				reply := &agentv1.StreamResponse{Principals: principals}
				if err := client.Stream.Send(reply); err != nil {
					slog.Error("updating agent error", "err", err)
				}
			}
		}

		return nil
	})

	s.PB.OnRecordBeforeDeleteRequest("machines").
		Add(func(e *core.RecordDeleteEvent) error {
			client := s.Clients[e.Record.Id]

			// Delete agent
			if err := client.Stream.Send(
				&agentv1.StreamResponse{Restore: BoolPointer(true)},
			); err != nil {
				slog.Error("deleting agent error", "err", err)
			}
			return nil
		})

	s.PB.OnRecordBeforeDeleteRequest("users").Add(func(e *core.RecordDeleteEvent) error {
		machines, err := getUserMachines(s.PB, e.Record)
		if err != nil {
			return err
		}
		for _, machine := range machines {
			principals, err := getPrincipals(s.PB, machine)
			if err != nil {
				return fmt.Errorf("failed to get machine users: %v", err)
			}

			client := s.Clients[machine.Id]
			reply := &agentv1.StreamResponse{Principals: principals}
			if err := client.Stream.Send(reply); err != nil {
				slog.Error("updating agent error", "err", err)
			}
		}
		return nil
	})
}
