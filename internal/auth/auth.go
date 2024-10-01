package auth

import (
	"github.com/llmariner/cli/internal/auth/apikeys"
	"github.com/spf13/cobra"
)

// Cmd is the root command for auth.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "auth",
		Short:              "Auth commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(apikeys.Cmd())
	cmd.AddCommand(loginCmd())

	cmd.AddCommand(statusCmd())
	return cmd
}
