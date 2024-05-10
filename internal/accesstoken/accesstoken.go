package accesstoken

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/llm-operator/cli/internal/config"
	"github.com/zchee/go-xdgbasedir"
	"gopkg.in/yaml.v2"
)

// T is the token.
type T struct {
	TokenType    string    `yaml:"tokenType"`
	TokenExpiry  time.Time `yaml:"tokenExpiry"`
	AccessToken  string    `yaml:"accessToken"`
	RefreshToken string    `yaml:"refreshToken"`
}

// saveToken saves the token to a file.
func saveToken(token *T) error {
	path := TokenFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config directory: %s", err)
	}
	b, err := yaml.Marshal(token)
	if err != nil {
		return fmt.Errorf("marshal token: %s", err)
	}
	if err := os.WriteFile(path, b, 0600); err != nil {
		return fmt.Errorf("write token: %s", err)
	}

	return nil
}

// LoadToken loads the token from a file.
func LoadToken(ctx context.Context, c *config.C) (*T, error) {
	path := TokenFilePath()
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read token: %s", err)
	}
	var token T
	if err := yaml.Unmarshal(b, &token); err != nil {
		return nil, fmt.Errorf("unmarshal token: %s", err)
	}

	tokenExchanger, err := NewTokenExchanger(c)
	if err != nil {
		return nil, fmt.Errorf("new token exchanger: %s", err)
	}
	token, err = tokenExchanger.refreshTokenIfExpired(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("refresh token: %s", err)
	}

	return &token, nil
}

// TokenFilePath returns the path to the token file.
func TokenFilePath() string {
	return filepath.Join(xdgbasedir.ConfigHome(), "llmo", "token.yaml")
}
