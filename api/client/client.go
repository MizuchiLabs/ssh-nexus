// Package client implements the gRPC server and client
package client

import (
	"context"
	"log/slog"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/MizuchiLabs/ssh-nexus/api/proto/gen/agent/v1/agentv1connect"
)

func Client(addr string) {
	for {
		conn, err := LoadCredentials(addr)
		if err != nil {
			slog.Error("failed to load credentials", "err", err)
			time.Sleep(3 * time.Second)
			continue
		}
		defer conn.CloseIdleConnections()

		client := agentv1connect.NewAgentServiceClient(conn, addr, connect.WithGRPC())

		// Authenticate
		token, err := GetToken()
		if err != nil {
			slog.Error("failed to get token", "err", err)
			return
		}

		agentID, err := GetAgentID()
		if err != nil {
			slog.Error("failed to get machine id", "err", err)
			return
		}

		hostname, err := os.Hostname()
		if err != nil {
			slog.Error("failed to get hostname", "err", err)
			hostname = "unknown"
		}

		ctx, cancel := context.WithCancel(context.Background())
		stream := client.Stream(ctx)
		stream.RequestHeader().Set("Authorization", "Bearer "+string(token))
		stream.RequestHeader().Set("AgentID", string(agentID))
		stream.RequestHeader().Set("Hostname", hostname)

		listener(ctx, stream)
		cancel()

		time.Sleep(3 * time.Second)
	}
}
