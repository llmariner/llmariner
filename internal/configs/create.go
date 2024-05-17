package configs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/llm-operator/cli/internal/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	configVersion    = "v1"
	authClientID     = "llm-operator"
	authClientSecret = "ZXhhbXBsZS1hcHAtc2VjcmV0"
	authRedirectURI  = "http://127.0.0.1:5555/callback"

	defaultEndpointURL = "http://localhost:8080/v1"
	defaultIssuerURL   = "http://kong-proxy.kong/v1/dex"
)

func createCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create()
		},
	}
}

func create() error {
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

	issuerURL, err := askWithDefaultValue(
		p,
		"Input OIDC Issuer URL for login",
		defaultIssuerURL,
		func(v string) error {
			if !isHTTPURL(v) {
				return errors.New("must start with 'http://' or 'https://'")
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	config := C{
		Version:     configVersion,
		EndpointURL: endpointURL,
		Auth: Auth{
			ClientID:     authClientID,
			ClientSecret: authClientSecret,
			RedirectURI:  authRedirectURI,
			IssuerURL:    issuerURL,
		},
	}

	path := configFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir all: %s", err)
	}
	b, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("marshal: %s", err)
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("write file: %s", err)
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
