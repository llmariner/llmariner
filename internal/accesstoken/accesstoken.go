package accesstoken

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zchee/go-xdgbasedir"
	"gopkg.in/yaml.v2"
)

// Token is the token.
type Token struct {
	AccessToken  string `yaml:"refreshToken"`
	RefreshToken string `yaml:"accessToken"`
}

// SaveToken saves the token to a file.
func SaveToken(token Token) error {
	path := filepath.Join(xdgbasedir.ConfigHome(), "llmo", "token.yaml")
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
