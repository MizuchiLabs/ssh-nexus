// Package config contains various utility functions and loads the configuration via envs
package config

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/MizuchiLabs/ssh-nexus/tools/util"
	"github.com/caarlos0/env/v11"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/settings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Config struct {
	// General Backend Settings
	AppURL           string `env:"PB_APP_URL"                 envDefault:""`
	AppName          string `env:"PB_APP_NAME"                envDefault:"SSH Nexus"`
	LogsMaxDays      int    `env:"PB_LOG_MAX_DAYS"            envDefault:"30"`
	AdminEmail       string `env:"PB_ADMIN_EMAIL"             envDefault:"root@nexus.local"`
	AdminPassword    string `env:"PB_ADMIN_PASSWORD,notEmpty" envDefault:""`
	SenderName       string `env:"PB_SENDER_NAME"             envDefault:"SSH Nexus"`
	SenderAddress    string `env:"PB_SENDER_EMAIL"            envDefault:"no-reply@nexus.local"`
	SMTPEnabled      bool   `env:"PB_SMTP_ENABLED"            envDefault:"false"`
	SMTPHost         string `env:"PB_SMTP_HOST"               envDefault:""`
	SMTPPort         int    `env:"PB_SMTP_PORT"               envDefault:"587"`
	SMTPUsername     string `env:"PB_SMTP_USER"               envDefault:""`
	SMTPPassword     string `env:"PB_SMTP_PASSWORD"           envDefault:""`
	SMTPTLS          bool   `env:"PB_SMTP_TLS"                envDefault:"false"`
	S3Enabled        bool   `env:"PB_S3_ENABLED"              envDefault:"false"`
	S3Endpoint       string `env:"PB_S3_ENDPOINT"             envDefault:""`
	S3Region         string `env:"PB_S3_REGION"               envDefault:""`
	S3Bucket         string `env:"PB_S3_BUCKET"               envDefault:"ssh-nexus"`
	S3Secret         string `env:"PB_S3_SECRET"               envDefault:""`
	S3AccessKey      string `env:"PB_S3_ACCESS_KEY"           envDefault:""`
	S3ForcePathStyle bool   `env:"PB_S3_FORCE_PATH_STYLE"     envDefault:"false"`
	OIDCURL          string `env:"PB_OIDC_URL"                envDefault:""`
	OIDCName         string `env:"PB_OIDC_NAME"               envDefault:""`
	OIDCRealm        string `env:"PB_OIDC_REALM"              envDefault:"master"`
	OIDCClientID     string `env:"PB_OIDC_CLIENT_ID"          envDefault:"ssh-nexus"`
	OIDCClientSecret string `env:"PB_OIDC_CLIENT_SECRET"      envDefault:""`

	// SSH Specific Settings
	DefaultRetention string `env:"DEFAULT_RETENTION" envDefault:"2592000"`
	UserLease        string `env:"USER_LEASE"        envDefault:"86400"`
	HostLease        string `env:"HOST_LEASE"        envDefault:"2592000"`
	MaxLease         string `env:"MAX_LEASE"         envDefault:"7776000"`
	SSHConfig        string `env:"SSH_CONFIG"        envDefault:""`
	InstallAgent     string `env:"INSTALL_AGENT"     envDefault:"true"`
}

// GetConfig returns the config used by the settings collection
func GetConfig() (*Config, error) {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("failed to parse envs: %v", err)
	}

	return &config, nil
}

// UseHTTPS returns true if the app url is a domain
func UseHTTPS() bool {
	config, err := GetConfig()
	if err != nil {
		return false
	}

	addr := net.ParseIP(config.AppURL)
	return addr == nil
}

// UpdateSettings sets various pocketbase and app related settings
func UpdateSettings(app core.App) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	// add admin if none found
	admin, err := app.Dao().FindAdminByEmail(config.AdminEmail)
	if err != nil {
		slog.Info("Setting up admin account...", "admin", config.AdminEmail)
		admin = &models.Admin{}
		admin.Email = config.AdminEmail
		if err = admin.SetPassword(config.AdminPassword); err != nil {
			return err
		}

		if err = app.Dao().SaveAdmin(admin); err != nil {
			return err
		}
	}

	// change admin password via env
	if ok := admin.ValidatePassword(config.AdminPassword); !ok {
		slog.Warn("password changed, setting new password for", "admin", config.AdminEmail)
		if err = admin.SetPassword(config.AdminPassword); err != nil {
			return err
		}
		if err = app.Dao().SaveAdmin(admin); err != nil {
			return err
		}
	}

	s, err := app.Dao().FindSettings(os.Getenv("PB_ENCRYPTION_KEY"))
	if err != nil {
		return err
	}

	if config.AppURL == "" {
		localIP := util.GetOutboundIP()
		s.Meta.AppUrl = fmt.Sprintf("http://%s:8090", localIP)
	} else {
		s.Meta.AppUrl = config.AppURL
	}
	s.Meta.AppName = config.AppName
	s.Logs.MaxDays = config.LogsMaxDays
	s.Meta.SenderName = config.SenderName
	s.Meta.SenderAddress = config.SenderAddress
	s.Smtp.Enabled = config.SMTPEnabled
	s.Smtp.Host = config.SMTPHost
	s.Smtp.Port = config.SMTPPort
	s.Smtp.Username = config.SMTPUsername
	s.Smtp.Password = config.SMTPPassword
	s.Smtp.Tls = config.SMTPTLS
	s.S3.Enabled = config.S3Enabled
	s.S3.Endpoint = config.S3Endpoint
	s.S3.Region = config.S3Region
	s.S3.Bucket = config.S3Bucket
	s.S3.Secret = config.S3Secret
	s.S3.AccessKey = config.S3AccessKey
	s.S3.ForcePathStyle = config.S3ForcePathStyle

	// OIDC helper
	if config.OIDCName != "" && config.OIDCClientID != "" && config.OIDCClientSecret != "" {
		oidc := OIDCConfig{
			URL:   config.OIDCURL,
			Name:  config.OIDCName,
			Realm: config.OIDCRealm,
		}
		if err = oidc.NewOIDC(); err != nil {
			slog.Error("failed to init oidc", "err", err)
		} else {
			s.OIDCAuth.Enabled = true
			s.OIDCAuth.DisplayName = cases.Title(language.English).String(config.OIDCName)
			s.OIDCAuth.ClientId = config.OIDCClientID
			s.OIDCAuth.ClientSecret = config.OIDCClientSecret
			s.OIDCAuth.AuthUrl = oidc.AuthEndpoint
			s.OIDCAuth.TokenUrl = oidc.TokenEndpoint
			s.OIDCAuth.UserApiUrl = oidc.UserInfoEndpoint
		}
	}

	// General settings (mostly ssh) to extend the app
	settingsColl, err := app.Dao().FindCollectionByNameOrId("settings")
	if err != nil {
		return err
	}

	if config.SSHConfig == "" {
		config.SSHConfig = `PubkeyAuthentication yes
PermitRootLogin yes
PermitEmptyPasswords no
PasswordAuthentication no
TrustedUserCAKeys /etc/ssh/nexus_user.pub
HostKey /etc/ssh/ssh_host_ed25519_key
HostCertificate /etc/ssh/ssh_host_ed25519_key-cert.pub
AuthorizedPrincipalsFile /etc/ssh/nexus_principals/%u`
	}

	baseSettings := map[string]string{
		"default_retention": config.DefaultRetention,
		"user_lease":        config.UserLease,
		"host_lease":        config.HostLease,
		"max_lease":         config.MaxLease,
		"install_agent":     config.InstallAgent,
		"ssh_config":        config.SSHConfig,
	}
	for k, v := range baseSettings {
		setting := models.NewRecord(settingsColl)
		setting.Set("key", k)
		setting.Set("value", v)
		if err = app.Dao().SaveRecord(setting); err != nil {
			continue
		}
	}

	// Email Templates
	s.Meta.VerificationTemplate = settings.EmailTemplate{
		Hidden:    false,
		Subject:   "Verify your email",
		ActionUrl: "{APP_URL}/_/#/auth/confirm-verification/{TOKEN}",
		Body: `<p>Hello,</p>
<p>Click on the button below to verify your email address.</p>
<p>
  <a class="btn" href="{ACTION_URL}" target="_blank" rel="noopener">Verify</a>
</p>
<p>
  Thanks,<br/>
  {APP_NAME} team
</p>`,
	}
	s.Meta.ResetPasswordTemplate = settings.EmailTemplate{
		Hidden:    false,
		Subject:   "Reset your password",
		ActionUrl: "{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}",
		Body: `<p>Hello,</p>
<p>Click on the button below to reset your password.</p>
<p>
  <a class="btn" href="{ACTION_URL}" target="_blank" rel="noopener">Reset password</a>
</p>
<p><i>If you didn't ask to reset your password, you can ignore this email.</i></p>
<p>
  Thanks,<br/>
  {APP_NAME} team
</p>`,
	}
	s.Meta.ConfirmEmailChangeTemplate = settings.EmailTemplate{
		Hidden:    false,
		Subject:   "Confirm your new email address",
		ActionUrl: "{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}",
		Body: `<p>Hello,</p>
<p>Click on the button below to confirm your new email address.</p>
<p>
  <a class="btn" href="{ACTION_URL}" target="_blank" rel="noopener">Confirm new email</a>
</p>
<p><i>If you didn't ask to change your email address, you can ignore this email.</i></p>
<p>
  Thanks,<br/>
  {APP_NAME} team
</p>`,
	}

	return app.Dao().SaveSettings(s)
}
