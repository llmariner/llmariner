package accesstoken

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/llmariner/llmariner/cli/internal/configs"
	"github.com/zchee/go-xdgbasedir"
	"golang.org/x/oauth2"
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
func saveToken(token *oauth2.Token) error {
	t := &T{
		TokenType:    token.TokenType,
		TokenExpiry:  token.Expiry,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	path := TokenFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config directory: %s", err)
	}
	b, err := yaml.Marshal(t)
	if err != nil {
		return fmt.Errorf("marshal token: %s", err)
	}
	if err := os.WriteFile(path, b, 0600); err != nil {
		return fmt.Errorf("write token: %s", err)
	}

	return nil
}

// LoadToken loads the token from a file.
func LoadToken(ctx context.Context, c *configs.C) (*T, error) {
	b, err := os.ReadFile(TokenFilePath())
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read token: %s", err)
		}
		// Fall back to the deprecated token file path.
		b, err = os.ReadFile(deprecatedTokenFilePath())
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("read token: %s", err)
			}
			return nil, fmt.Errorf("token file not found. Please run 'llma auth login'")
		}
	}

	var token T
	if err := yaml.Unmarshal(b, &token); err != nil {
		return nil, fmt.Errorf("unmarshal token: %s", err)
	}

	token, err = refreshTokenIfExpired(ctx, token, c.Auth)
	if err != nil {
		return nil, fmt.Errorf("refresh token: %s", err)
	}

	return &token, nil
}

// TokenFilePath returns the path to the token file.
func TokenFilePath() string {
	return filepath.Join(xdgbasedir.ConfigHome(), "llmariner", "token.yaml")
}

func deprecatedTokenFilePath() string {
	return filepath.Join(xdgbasedir.ConfigHome(), "llmo", "token.yaml")
}
