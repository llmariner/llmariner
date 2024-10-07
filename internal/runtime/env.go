package runtime

import (
	"context"
	"fmt"

	"github.com/llmariner/llmariner/internal/accesstoken"
	"github.com/llmariner/llmariner/internal/configs"
)

// NewEnv creates a new runtime env.
func NewEnv(ctx context.Context) (*Env, error) {
	c, err := configs.LoadOrCreate()
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
	Config *configs.C
	Token  *accesstoken.T
}
