package auth

import (
	"context"
	"net/http"

	"github.com/llmariner/llmariner/cli/internal/accesstoken"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/llmariner/llmariner/cli/internal/ui"
	uv1 "github.com/llmariner/user-manager/api/v1"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Login status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return status(cmd.Context())
		},
	}
}

func status(ctx context.Context) error {
	p := ui.NewPrompter()

	env, err := runtime.NewEnv(ctx)
	if err != nil {
		p.Warn("Not logged in.")
		return nil
	}

	var req uv1.GetUserSelfRequest
	var user uv1.User
	if err := ihttp.NewClient(env).Send(http.MethodGet, "/users:getSelf", &req, &user); err != nil {
		return err
	}

	p.Printf("Logged in.\n")
	p.Printf("- ID: %s\n", user.Id)
	p.Printf("- Token file location: %s\n", accesstoken.TokenFilePath())

	return nil
}
