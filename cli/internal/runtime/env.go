package runtime

import (
	"context"
	"fmt"

	"github.com/llmariner/llmariner/cli/internal/accesstoken"
	"github.com/llmariner/llmariner/cli/internal/configs"
)

// NewEnv creates a new runtime env.
func NewEnv(ctx context.Context) (*Env, error) {
	c, err := configs.LoadOrCreate()
	if err != nil {
		return nil, fmt.Errorf("load or create config: %s", err)
	}

	// If the API key is set in the env var, use it; we don't need to load the token.
	if v := accesstoken.GetAPIKeyEnvVar(); v != "" {
		return &Env{
			Config:       c,
			APIKeySecret: v,
		}, nil
	}

	t, err := accesstoken.LoadToken(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("load token: %s", err)
	}

	return &Env{
		Config: c,
		Token:  t,
	}, nil
}

// Env is a struct that contains the runtime env for the CLI.
type Env struct {
	Config       *configs.C
	APIKeySecret string
	Token        *accesstoken.T
}

// AccessToken returns the access token.
func (e *Env) AccessToken() string {
	if e.APIKeySecret != "" {
		return e.APIKeySecret
	}
	return e.Token.AccessToken
}
