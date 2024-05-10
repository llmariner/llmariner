package runtime

import (
	"context"
	"fmt"

	"github.com/llm-operator/cli/internal/accesstoken"
	"github.com/llm-operator/cli/internal/config"
)

// NewEnv creates a new runtime env.
func NewEnv(ctx context.Context) (*Env, error) {
	c, err := config.LoadOrCreate()
	if err != nil {
		return nil, fmt.Errorf("load or create config: %s", err)
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
	Config *config.C
	Token  *accesstoken.T
}
