package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zchee/go-xdgbasedir"
	"gopkg.in/yaml.v2"
)

var (
	defaultConfig = C{
		Version: "v1",
		Auth: Auth{
			ClientID:     "llm-operator",
			ClientSecret: "ZXhhbXBsZS1hcHAtc2VjcmV0",
			RedirectURI:  "http://127.0.0.1:5555/callback",
			IssuerURL:    "http://kong-kong-proxy.kong/v1/dex",
		},
	}
)

// Auth is an authentication configuration.
type Auth struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	RedirectURI  string `yaml:"redirectUri"`
	IssuerURL    string `yaml:"issuerUrl"`
}

func (a *Auth) validate() error {
	if a.ClientID == "" {
		return errors.New("clientId is required")
	}
	if a.ClientSecret == "" {
		return errors.New("clientSecret is required")
	}
	if a.RedirectURI == "" {
		return errors.New("redirectUri is required")
	}
	if a.IssuerURL == "" {
		return errors.New("issuerUrl is required")
	}
	return nil
}

// C is a config file.
type C struct {
	Version string `yaml:"version"`
	Auth    Auth   `yaml:"auth"`
}

// Validate validates the config.
func (c *C) Validate() error {
	if c.Version == "" {
		return errors.New("version is required")
	}

	if err := c.Auth.validate(); err != nil {
		return fmt.Errorf("auth: %s", err)
	}
	return nil
}

// Load loads the config.
func Load() (*C, error) {
	path := filepath.Join(xdgbasedir.ConfigHome(), "llmo", "config.yaml")

	// Create a config file if it doesn't exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := createDefaultConfig(path); err != nil {
			return nil, fmt.Errorf("create default config: %s", err)
		}
	}

	var config C
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %s", err)
	}

	if err = yaml.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf("unmarshal: %s", err)
	}
	return &config, nil
}

func createDefaultConfig(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir all: %s", err)
	}
	// Create a new config file.
	b, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("marshal: %s", err)
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("write file: %s", err)
	}
	return nil

}
