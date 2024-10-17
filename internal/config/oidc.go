package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OIDCConfig struct {
	Name             string `json:"name,omitempty"`
	URL              string `json:"url,omitempty"`
	Realm            string `json:"realm,omitempty"`
	AuthEndpoint     string `json:"authorization_endpoint,omitempty"`
	TokenEndpoint    string `json:"token_endpoint,omitempty"`
	UserInfoEndpoint string `json:"userinfo_endpoint,omitempty"`
}

// NewOIDC creates a new provider
func (o *OIDCConfig) NewOIDC() error {
	switch strings.ToLower(o.Name) {
	case "keycloak":
		return o.getKeycloakConfig()
	case "vault":
		// TODO
		return nil
	default:
		return fmt.Errorf("unsupported oidc provider: %s", o.Name)
	}
}

func (o *OIDCConfig) getKeycloakConfig() error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/realms/%s/.well-known/openid-configuration", o.URL, o.Realm),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("(%d) failed to send oidc config request", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&o); err != nil {
		return err
	}

	return nil
}
