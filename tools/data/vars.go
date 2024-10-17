package data

import (
	"os"
	"path/filepath"
)

var (
	// Certificates & Keys for the grpc server
	BaseCertDir  = Path("certs")
	ServerCaCert = Path("certs/proto_server_ca.pem")
	ServerCaKey  = Path("certs/proto_server_ca_key.pem")
	ServerCert   = Path("certs/proto_server.pem")
	ServerKey    = Path("certs/proto_server_key.pem")

	// User key to verify users against the host
	UserKey = Path("nexus_user.key")

	// CA to verify hosts against the user
	HostCAKey = Path("nexus_host_ca.key")

	// Secret to authenticate the agent
	Token = Path("token")

	// Various paths used on the server
	SSHConfigPath      = "/etc/ssh/sshd_config.d/nexus.conf"
	PrincipalPath      = "/etc/ssh/nexus_principals/"
	PublicUserKeyPath  = "/etc/ssh/nexus_user.pub"
	PrivateHostKeyPath = "/etc/ssh/ssh_host_ed25519_key"
	PublicHostKeyPath  = "/etc/ssh/ssh_host_ed25519_key.pub"
	CertHostPath       = "/etc/ssh/ssh_host_ed25519_key-cert.pub"
	AuthorizedKeysPath = "~/.ssh/authorized_keys"

	// Path to the agent binary
	AgentPath    = "/usr/local/bin/nexus-agent"
	AgentService = "/etc/systemd/system/nexus-agent.service"

	// Path to the temporary agent binary (for downloads and updates)
	AgentDownloadPath = filepath.Join(os.TempDir(), "nexus-agent")
)
