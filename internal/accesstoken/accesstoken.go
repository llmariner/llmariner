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
	AccessToken  string `yaml:"accessToken"`
	RefreshToken string `yaml:"refreshToken"`
}

// SaveToken saves the token to a file.
func SaveToken(token *T) error {
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
func LoadToken() (*T, error) {
	path := TokenFilePath()
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

// TokenFilePath returns the path to the token file.
func TokenFilePath() string {
	return filepath.Join(xdgbasedir.ConfigHome(), "llmo", "token.yaml")
}
