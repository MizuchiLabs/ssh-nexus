package client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"connectrpc.com/connect"
	agentv1 "github.com/MizuchiLabs/ssh-nexus/api/proto/gen/agent/v1"
	"github.com/MizuchiLabs/ssh-nexus/tools/data"
	"github.com/MizuchiLabs/ssh-nexus/tools/updater"
	"golang.org/x/crypto/ssh"
)

// listener sends a request to the server and listens for responses
func listener(
	ctx context.Context,
	stream *connect.BidiStreamForClient[agentv1.StreamRequest, agentv1.StreamResponse],
) {
	if err := stream.Send(createRequest()); err != nil {
		slog.Error("failed to send request", "err", err)
		return
	}

	go monitorCertificate(ctx, stream)

	for {
		resp, err := stream.Receive()
		if err != nil {
			slog.Error(
				"server disconnected, reconnecting...",
				"err",
				err.Error(),
			)
			return
		}
		if err := action(resp); err != nil {
			slog.Error("failed to update files", "err", err)
		}
	}
}

func createRequest() *agentv1.StreamRequest {
	request := agentv1.StreamRequest{}

	pubHostKey, err := getPublicHostKey()
	if err != nil {
		slog.Error("failed to get host key", "err", err)
	}

	request.Version = &updater.Version
	request.PublicHostKey = &pubHostKey

	return &request
}

// monitorCertificate checks if the host certificate needs to be renewed
func monitorCertificate(
	ctx context.Context,
	stream *connect.BidiStreamForClient[agentv1.StreamRequest, agentv1.StreamResponse],
) {
	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			renew, err := renewHostCert()
			if err != nil {
				slog.Error("failed to check certificate renewal", "err", err)
				continue
			}
			if renew {
				if err := stream.Send(createRequest()); err != nil {
					slog.Error("failed to send request", "err", err)
					return
				}
			}
			<-ticker.C
		}
	}
}

func renewHostCert() (bool, error) {
	if _, err := os.Stat(data.CertHostPath); os.IsNotExist(err) {
		return true, nil
	}

	hostCert, err := os.ReadFile(data.CertHostPath)
	if err != nil {
		return false, fmt.Errorf("failed to read host cert: %w", err)
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(hostCert)
	if err != nil {
		return false, fmt.Errorf("failed to parse host cert: %w", err)
	}

	cert, ok := pubKey.(*ssh.Certificate)
	if !ok {
		return false, fmt.Errorf("not a valid certificate: %w", err)
	}

	// Renew 10 days before the certificate expires
	if time.Now().UTC().AddDate(0, 0, 10).After(time.Unix(int64(cert.ValidBefore), 0).UTC()) {
		slog.Info("Certificate will be renewed", "expiry", cert.ValidBefore)
		return true, nil
	}
	return false, nil
}

// getPublicHostKey checks if a default host key exists, creates it if not and returns the public host key
func getPublicHostKey() (string, error) {
	if _, err := os.Stat(data.PrivateHostKeyPath); os.IsNotExist(err) {
		if err := data.NewSigner(data.PrivateHostKeyPath, "host@ssh-nexus"); err != nil {
			return "", fmt.Errorf("failed to create host ca: %w", err)
		}
	}

	pubHostKey, err := os.ReadFile(data.PublicHostKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read host private key: %w", err)
	}

	return strings.TrimSpace(string(pubHostKey)), nil
}

func action(resp *agentv1.StreamResponse) error {
	if err := updateSSHConfig(resp.GetSshConfig()); err != nil {
		return err
	}
	if err := updateUserCA(resp.GetUserCertificatePublicKey()); err != nil {
		return err
	}
	if err := updateHostCert(resp.GetHostCertificatePublicKey()); err != nil {
		return err
	}
	if err := updatePrincipals(resp.GetPrincipals()); err != nil {
		return err
	}

	if resp.GetRestore() {
		if err := restore(); err != nil {
			return err
		}
	}
	return nil
}

// Add our custom ssh config to the server
func updateSSHConfig(config []byte) error {
	if config == nil {
		return nil
	}
	slog.Info("updated ssh config")
	return os.WriteFile(data.SSHConfigPath, config, 0600)
}

// Add our public key to the server (as a user ca and authorized key)
func updateUserCA(pub []byte) error {
	if pub == nil {
		return nil
	}
	slog.Info("updated user ca")
	if err := os.WriteFile(data.AuthorizedKeysPath, pub, 0600); err != nil {
		return fmt.Errorf("failed to write authorized keys: %w", err)
	}
	return os.WriteFile(data.PublicUserKeyPath, pub, 0600)
}

// Add a custom private host key for identification
func updateHostCert(pub []byte) error {
	if pub == nil {
		return nil
	}
	slog.Info("updated host certificate")
	return os.WriteFile(data.CertHostPath, pub, 0600)
}

// Add correct principals to the server
func updatePrincipals(principals []*agentv1.StreamResponse_Principal) error {
	principalMap := make(map[string][]string)
	for _, principal := range principals {
		principalMap[principal.GetKey()] = append(
			principalMap[principal.GetKey()],
			principal.GetValues()...)
	}
	if !slices.Contains(principalMap["root"], "root") {
		principalMap["root"] = append([]string{"root"}, principalMap["root"]...)
	}

	groups := make(map[string]bool) // Used to check if the groups are the same
	for group, users := range principalMap {
		GroupPath := filepath.Join(data.PrincipalPath, group)

		groups[group] = true

		if err := os.MkdirAll(data.PrincipalPath, 0750); err != nil {
			return err
		}

		if err := os.WriteFile(GroupPath, []byte(strings.Join(users, "\n")), 0600); err != nil {
			return err
		}
		slog.Info("updated principals", "group", group, "users", users)
	}

	// Get the current groups on the machine and compare them
	currentGroups, err := os.ReadDir(data.PrincipalPath)
	if err != nil {
		return err
	}

	// Check if the groups are the same if not, delete them
	for _, group := range currentGroups {
		if !groups[group.Name()] {
			GroupPath := filepath.Join(data.PrincipalPath, group.Name())
			if err := os.RemoveAll(GroupPath); err != nil {
				return err
			}
			slog.Info("deleted principals", "group", group.Name())
		}
	}

	return nil
}

// restore removes all files from the server for cleanup
func restore() error {
	dirs := []string{
		data.PrincipalPath,
		data.SSHConfigPath,
		data.PublicUserKeyPath,
		data.CertHostPath,
		data.AgentPath,
		data.AgentService,
		data.Token,
	}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}

	if err := exec.Command("systemctl", "disable", "nexus-agent").Run(); err != nil {
		return err
	}
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}
