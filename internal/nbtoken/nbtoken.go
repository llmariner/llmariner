package nbtoken

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zchee/go-xdgbasedir"
	"gopkg.in/yaml.v2"
)

const version = "v1"

type notebookTokens struct {
	Version string            `yaml:"version"`
	Tokens  map[string]string `yaml:"tokens"`
}

// SaveToken saves the token to a file.
func SaveToken(nbID, token string) error {
	t, err := loadTokens()
	if err != nil {
		return err
	}
	t.Tokens[nbID] = token

	path := tokenFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	fileData, err := yaml.Marshal(&t)
	if err != nil {
		return err
	}
	return os.WriteFile(path, fileData, 0600)
}

// LoadToken loads the token from a file.
func LoadToken(nbID string) (string, error) {
	t, err := loadTokens()
	if err != nil {
		return "", err
	}

	token, ok := t.Tokens[nbID]
	if !ok {
		return "", fmt.Errorf("notebook token not found")
	}
	return token, nil
}

func loadTokens() (*notebookTokens, error) {
	file, err := os.ReadFile(tokenFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return &notebookTokens{
				Version: version,
				Tokens:  make(map[string]string),
			}, nil
		}
		return nil, err
	}

	var t notebookTokens
	if err := yaml.Unmarshal(file, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func tokenFilePath() string {
	return filepath.Join(xdgbasedir.ConfigHome(), "llmo", "notebook_tokens.yaml")
}
