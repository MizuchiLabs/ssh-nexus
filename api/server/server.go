// Package server implements the gRPC server and client
package server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	agentv1 "github.com/MizuchiLabs/ssh-nexus/api/proto/gen/agent/v1"
	"github.com/MizuchiLabs/ssh-nexus/api/proto/gen/agent/v1/agentv1connect"
	"github.com/MizuchiLabs/ssh-nexus/tools/data"
	"github.com/MizuchiLabs/ssh-nexus/tools/updater"
	"github.com/MizuchiLabs/ssh-nexus/tools/util"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/crypto/acme/autocert"
)

type AgentServer struct {
	agentv1connect.UnimplementedAgentServiceHandler
	mu      sync.Mutex
	PB      core.App
	Clients map[string]Client
}

type Client struct {
	Machine *models.Record
	Stream  *connect.BidiStream[agentv1.StreamRequest, agentv1.StreamResponse]
}

func (s *AgentServer) Stream(
	ctx context.Context,
	stream *connect.BidiStream[agentv1.StreamRequest, agentv1.StreamResponse],
) error {
	clientID, err := s.connect(stream)
	if err != nil {
		return err
	}

	s.mu.Lock()
	client := s.Clients[*clientID]
	s.mu.Unlock()

	s.pbHook()
	defer s.disconnect(client)

	for {
		if err := ctx.Err(); err != nil {
			return connect.NewError(connect.CodeUnknown, err)
		}

		request, err := stream.Receive()
		if err != nil {
			// Stream closed by client
			if errors.Is(err, io.EOF) {
				break
			}
			slog.Error("failed to receive", "err", err)
			return connect.NewError(connect.CodeUnknown, err)
		}

		s.monitorHook(client, request)
	}

	return nil
}

func (s *AgentServer) connect(
	stream *connect.BidiStream[agentv1.StreamRequest, agentv1.StreamResponse],
) (*string, error) {
	if err := validate(stream.RequestHeader()); err != nil {
		return nil, err
	}
	agentID := stream.RequestHeader().Get("AgentID")
	if agentID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("missing agent id"))
	}
	hostname := stream.RequestHeader().Get("Hostname")
	if hostname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("missing hostname"))
	}

	host, _, err := net.SplitHostPort(stream.Peer().Addr)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	machine, err := filterMachine(s.PB, host, agentID, hostname)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.mu.Lock()
	s.Clients[machine.Id] = Client{Machine: machine, Stream: stream}
	s.mu.Unlock()
	slog.Info("client connected", "id", machine.Id, "name", machine.GetString("name"))

	s.initClient(s.Clients[machine.Id])

	return &machine.Id, nil
}

func (s *AgentServer) disconnect(client Client) {
	s.mu.Lock()
	delete(s.Clients, client.Machine.Id)
	s.mu.Unlock()

	slog.Info(
		"client disconnected",
		"id",
		client.Machine.Id,
		"name",
		client.Machine.GetString("name"),
	)

	client.Machine.Set("agent", false)
	if err := s.PB.Dao().SaveRecord(client.Machine); err != nil {
		slog.Error("failed to save machine", "err", err)
	}
}

func validate(header http.Header) error {
	auth := header.Get("authorization")
	if len(auth) == 0 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("missing authorization"))
	}
	if !strings.HasPrefix(auth, "Bearer ") {
		return connect.NewError(connect.CodeUnauthenticated, errors.New("missing bearer prefix"))
	}

	token, err := os.ReadFile(data.Token)
	if err != nil {
		return connect.NewError(connect.CodeInternal, errors.New("failed to read token"))
	}

	if strings.TrimSpace(string(token)) != strings.TrimPrefix(auth, "Bearer ") {
		return connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("failed to validate password"),
		)
	}

	return err
}

func Server(app core.App) {
	agentServer := &AgentServer{
		PB:      app,
		Clients: make(map[string]Client),
	}

	settings, err := app.Dao().FindSettings(os.Getenv("PB_ENCRYPTION_KEY"))
	if err != nil {
		return
	}

	var tlsConfig *tls.Config
	if util.IsIP(settings.Meta.AppUrl) {
		certs, err := tls.LoadX509KeyPair(data.ServerCert, data.ServerKey)
		if err != nil {
			slog.Error("failed to load server cert", "err", err)
			return
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{certs},
			MinVersion:   tls.VersionTLS12,
		}
	} else {
		url, _ := url.Parse(settings.Meta.AppUrl)
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(url.Hostname()),
			Cache:      autocert.DirCache(data.BaseCertDir),
		}
		tlsConfig = certManager.TLSConfig()
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))

	}

	mux := http.NewServeMux()
	mux.Handle(agentv1connect.NewAgentServiceHandler(agentServer))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.RequestURI(), http.StatusMovedPermanently)
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(updater.Version)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/ca.crt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, data.ServerCaCert)
	})
	mux.HandleFunc("/client/{id}", func(w http.ResponseWriter, r *http.Request) {
		client := agentServer.Clients[r.PathValue("id")]
		if client.Machine == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		machine, err := client.Machine.MarshalJSON()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(machine); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("/client/{id}/health", func(w http.ResponseWriter, r *http.Request) {
		client := agentServer.Clients[r.PathValue("id")]
		if client.Machine == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if _, err := w.Write([]byte("OK")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	srv := &http.Server{
		Addr:              ":8091",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
		TLSConfig:         tlsConfig,
	}

	slog.Info("gRPC server running on", "port", ":8091")
	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil {
			if err == http.ErrServerClosed {
				slog.Info("gRPC server closed")
				return
			}
			slog.Error("gRPC server error", "err", err)
			return
		}
	}()
}
