package configs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/llmariner/llmariner/internal/ui"
	"github.com/zchee/go-xdgbasedir"
	"gopkg.in/yaml.v2"
)

const (
	configVersion    = "v1"
	authClientID     = "llmariner"
	authClientSecret = "ZXhhbXBsZS1hcHAtc2VjcmV0"
	authRedirectURI  = "http://127.0.0.1:5555/callback"

	defaultEndpointURL = "http://localhost:8080/v1"
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

// Context is a context configuration.
type Context struct {
	OrganizationID string `yaml:"organizationId"`
	ProjectID      string `yaml:"projectId"`
}

// C is a config file.
type C struct {
	Version string `yaml:"version"`

	EndpointURL string `yaml:"endpointUrl"`

	Auth Auth `yaml:"auth"`

	Context Context `yaml:"context"`
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

// Save saves the config.
func (c *C) Save() error {
	path := configFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir all: %s", err)
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal: %s", err)
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("write file: %s", err)
	}
	return nil
}

// LoadOrCreate loads the config.
func LoadOrCreate() (*C, error) {
	// Create a config file if it doesn't exists.
	path := configFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Check if the deprecated config file exists.
		deprecatedPath := deprecatedConfigFilePath()
		if _, err := os.Stat(deprecatedPath); err == nil {
			path = deprecatedPath
		} else {
			if err := CreateNewConfig(); err != nil {
				return nil, fmt.Errorf("create default config: %s", err)
			}
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

// CreateNewConfig creates a new config file.
func CreateNewConfig() error {
	p := ui.NewPrompter()
	endpointURL, err := askWithDefaultValue(
		p,
		"Input endpoint URL of LLM service",
		defaultEndpointURL,
		func(v string) error {
			if !isHTTPURL(v) {
				return errors.New("must start with 'http://' or 'https://'")
			}
			if !strings.HasSuffix(v, "/v1") {
				return errors.New("must end with '/v1'")
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	// Remove the trailing slash.
	endpointURL = strings.TrimSuffix(endpointURL, "/")

	issuerURL := fmt.Sprintf("%s/dex", endpointURL)
	c := &C{
		Version:     configVersion,
		EndpointURL: endpointURL,
		Auth: Auth{
			ClientID:     authClientID,
			ClientSecret: authClientSecret,
			RedirectURI:  authRedirectURI,
			IssuerURL:    issuerURL,
		},
	}
	if err := c.Save(); err != nil {
		return err
	}

	return nil
}

func askWithDefaultValue(
	p ui.Prompter,
	msg string,
	defaultVal string,
	validate func(v string) error,
) (string, error) {
	var val string
	if err := p.Ask(
		&survey.Input{Message: fmt.Sprintf("%s (default: %q):", msg, defaultVal)},
		&val,
		survey.WithValidator(func(ans interface{}) error {
			v, ok := ans.(string)
			if !ok {
				return fmt.Errorf("invalid input")
			}
			v = strings.TrimSpace(v)
			if v == "" {
				return nil
			}

			return validate(v)
		}),
	); err != nil {
		return "", err
	}
	if val != "" {
		return val, nil
	}

	return defaultVal, nil
}

func isHTTPURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func configFilePath() string {
	return filepath.Join(xdgbasedir.ConfigHome(), "llmariner", "config.yaml")
}

func deprecatedConfigFilePath() string {
	return filepath.Join(xdgbasedir.ConfigHome(), "llmo", "config.yaml")
}
