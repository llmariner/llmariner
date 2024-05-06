package auth

import (
	"github.com/llm-operator/cli/internal/auth/apikeys"
	"github.com/spf13/cobra"
)

// Cmd is the root command for auth.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "auth",
		Short:              "Auth commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: false,
	}
	cmd.AddCommand(apikeys.Cmd())
	cmd.AddCommand(loginCmd())
	return cmd
}
