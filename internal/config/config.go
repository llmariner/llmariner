package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zchee/go-xdgbasedir"
	"gopkg.in/yaml.v2"
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

	EndpointURL string `yaml:"endpointUrl"`

	Auth Auth `yaml:"auth"`
}

// Validate validates the config.
func (c *C) Validate() error {
	if c.Version == "" {
		return errors.New("version is required")
	}

	if c.EndpointURL == "" {
		return errors.New("endpointUrl is required")
	}

	if err := c.Auth.validate(); err != nil {
		return fmt.Errorf("auth: %s", err)
	}
	return nil
}

// LoadOrCreate loads the config.
func LoadOrCreate() (*C, error) {
	// Create a config file if it doesn't exists.
	path := configFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := create(); err != nil {
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

func configFilePath() string {
	return filepath.Join(xdgbasedir.ConfigHome(), "llmo", "config.yaml")
}

// Cmd returns a new config command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "config",
		Short:              "Config commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: false,
	}
	cmd.AddCommand(createCmd())
	return cmd
}
