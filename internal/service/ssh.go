package service

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MizuchiLabs/ssh-nexus/tools/data"
	"github.com/MizuchiLabs/ssh-nexus/tools/updater"
	"github.com/MizuchiLabs/ssh-nexus/tools/util"
	"github.com/pkg/sftp"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/crypto/ssh"
)

// ManualUpdate for machines running without an agent
func ManualUpdate(app core.App, machine *models.Record) {
	if machine.GetBool("agent") {
		slog.Debug("Skipping manual update for machine", "name", machine.GetString("name"))
		return
	}
	defer func() {
		if err := app.Dao().SaveRecord(machine); err != nil {
			slog.Error("Failed to save machine", "name", machine.GetString("name"), "err", err)
		}
	}()

	conn, err := connect(machine)
	if err != nil {
		machine.Set("error", err.Error())
		return
	}
	defer conn.Close()

	// Prepare SSH Configuration & Principals
	sshConfig, err := app.Dao().FindFirstRecordByData("settings", "key", "ssh_config")
	if err != nil {
		machine.Set("error", err.Error())
		return
	}
	publicKeyFile, err := data.GetPublicUserKey()
	if err != nil {
		machine.Set("error", err.Error())
		return
	}
	commands := []string{
		fmt.Sprintf("mkdir -p %s", data.PrincipalPath),
		fmt.Sprintf("echo -n '%s' | tee %s", string(publicKeyFile), data.PublicUserKeyPath),
		fmt.Sprintf("echo -n '%s' | tee %s", sshConfig.GetString("value"), data.SSHConfigPath),
	}
	command := strings.Join(commands, "; ")
	if _, err = run(conn, command); err != nil {
		machine.Set("error", err.Error())
		return
	}

	groups, err := GetMachineUsers(app, machine)
	if err != nil {
		machine.Set("error", err.Error())
		return
	}

	if err := setPrincipals(conn, groups); err != nil {
		machine.Set("error", err.Error())
		return
	}
}

// TODO: Check if agent is installed and running, if not try to restart
func InstallAgent(app core.App, machine *models.Record) {
	if machine.GetBool("agent") {
		slog.Debug("Skipping install for machine", "name", machine.GetString("name"))
		return
	}
	defer func() {
		if err := app.Dao().SaveRecord(machine); err != nil {
			slog.Error("Failed to save machine", "name", machine.GetString("name"), "err", err)
		}
	}()

	settings, err := app.Dao().FindSettings(os.Getenv("PB_ENCRYPTION_KEY"))
	if err != nil {
		machine.Set("error", err.Error())
		return
	}

	conn, err := connect(machine)
	if err != nil {
		machine.Set("error", err.Error())
		return
	}
	defer conn.Close()

	agentToken, err := data.GetToken()
	if err != nil {
		machine.Set("error", err.Error())
		return
	}

	if err = uploadAgent(conn); err != nil {
		machine.Set("error", err.Error())
		return
	}

	appURL, err := url.Parse(settings.Meta.AppUrl)
	if err != nil {
		machine.Set("error", err.Error())
		return
	}

	agentPath := "/usr/local/bin/nexus-agent"
	systemdPath := "/etc/systemd/system/nexus-agent.service"
	systemdConfig := fmt.Sprintf(`[Unit]
Description=Nexus Agent
After=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=%s -server %s
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target`, agentPath, appURL.Hostname())

	// Install Agent
	commands := []string{
		fmt.Sprintf("chmod +x %s", agentPath),
		fmt.Sprintf("echo -n '%s' | tee %s", systemdConfig, systemdPath),
		fmt.Sprintf("echo -n '%s' | tee %s", agentToken, data.Token),
		"systemctl daemon-reload",
		"systemctl enable nexus-agent",
		"systemctl start nexus-agent",
	}
	command := strings.Join(commands, "; ")
	if _, err := run(conn, command); err != nil {
		machine.Set("error", err.Error())
		return
	}
}

func Restore(machine *models.Record) {
	if machine.GetBool("agent") {
		slog.Debug("Skipping restore for machine", "name", machine.GetString("name"))
		return
	}

	conn, err := connect(machine)
	if err != nil {
		slog.Error("Failed to connect to machine", "name", machine.GetString("name"), "err", err)
		return
	}
	defer conn.Close()

	commands := []string{
		fmt.Sprintf("rm -rf %s", data.PrincipalPath),
		fmt.Sprintf("rm %s", data.SSHConfigPath),
		fmt.Sprintf("rm %s", data.PublicUserKeyPath),
		fmt.Sprintf("rm %s", data.PublicHostKeyPath),
		fmt.Sprintf("rm %s", data.AgentPath),
		fmt.Sprintf("rm %s", data.AgentService),
		fmt.Sprintf("rm %s", data.Token),
		"systemctl disable nexus-agent.service",
		"systemctl stop nexus-agent.service",
		"systemctl daemon-reload",
	}

	command := strings.Join(commands, "; ")
	if _, err := run(conn, command); err != nil {
		slog.Error("Failed to restore machine", "name", machine.GetString("name"), "err", err)
		return
	}
}

func uploadAgent(conn *ssh.Client) error {
	if conn == nil {
		return fmt.Errorf("no connection provided")
	}

	if err := updater.CheckAgent(); err != nil {
		return err
	}

	localAgent, err := os.Open(data.AgentDownloadPath)
	if err != nil {
		return err
	}
	defer localAgent.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	remoteAgentPath := "/usr/local/bin/nexus-agent"
	remoteAgent, err := client.Create(remoteAgentPath)
	if err != nil {
		return err
	}
	defer remoteAgent.Close()

	_, err = io.Copy(remoteAgent, localAgent)
	if err != nil {
		return err
	}

	if err := client.Chmod(remoteAgentPath, 0755); err != nil {
		return err
	}

	return nil
}

// UpdatePrincipal checks the principals on a machine and updates them if necessary.
func setPrincipals(conn *ssh.Client, groups map[string][]string) error {
	var groupsList []string
	for group, principals := range groups {
		groupsList = append(groupsList, group)
		GroupPath := filepath.Join(data.PrincipalPath, group)

		commands := []string{
			fmt.Sprintf("touch %s", GroupPath),
			fmt.Sprintf("echo -e '%s' | tee %s", strings.Join(principals, "\n"), GroupPath),
		}
		command := strings.Join(commands, "; ")
		if _, err := run(conn, command); err != nil {
			return err
		}
	}

	// Get the current groups on the machine and compare them
	currentGroups, err := run(conn, "ls -1 "+data.PrincipalPath)
	if err != nil {
		return err
	}

	// Check if the groups are the same if not, delete them
	groupState := strings.Split(strings.TrimSpace(string(currentGroups)), "\n")
	groupDiff := util.Diff(groupState, groupsList)
	if len(groupDiff) > 0 {
		for _, group := range groupDiff {
			if _, err = run(conn, "rm -f "+data.PrincipalPath+group); err != nil {
				return err
			}
		}
	}

	return nil
}

// Opens a new SSH connection
func connect(machine *models.Record) (*ssh.Client, error) {
	if machine == nil {
		return nil, fmt.Errorf("no machine provided")
	}

	signer, err := data.GetUserSigner()
	if err != nil {
		return nil, err
	}

	addr := net.JoinHostPort(machine.GetString("host"), machine.GetString("port"))
	conn, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func run(client *ssh.Client, command string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		session, err := client.NewSession()
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %v", err)
		}
		defer session.Close()

		type result struct {
			output []byte
			err    error
		}
		resultChan := make(chan result, 1)

		go func() {
			output, err := session.CombinedOutput(command)
			resultChan <- result{output: output, err: err}
		}()

		// Wait for command completion or timeout
		select {
		case res := <-resultChan:
			if res.err == nil {
				return res.output, nil
			}
			lastErr = res.err
		case <-time.After(10 * time.Second):
			lastErr = fmt.Errorf("command timed out")
		}

		// Close the current session to prevent resource leaks before retrying
		session.Close()

		// Wait for a bit before retrying
		time.Sleep(1 * time.Second)
	}
	return nil, fmt.Errorf("error: %v", lastErr)
}
