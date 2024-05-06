package accesstoken

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zchee/go-xdgbasedir"
	"gopkg.in/yaml.v2"
)

// T is the token.
type T struct {
	AccessToken  string `yaml:"refreshToken"`
	RefreshToken string `yaml:"accessToken"`
}

// SaveToken saves the token to a file.
func SaveToken(token *T) error {
	path := tokenFilePath()
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
func LoadToken() (*T, error) {
	path := tokenFilePath()
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read token: %s", err)
	}
	var token T
	if err := yaml.Unmarshal(b, &token); err != nil {
		return nil, fmt.Errorf("unmarshal token: %s", err)
	}
	return &token, nil
}

func tokenFilePath() string {
	return filepath.Join(xdgbasedir.ConfigHome(), "llmo", "token.yaml")
}
