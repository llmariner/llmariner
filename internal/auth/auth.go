package auth

import (
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
	cmd.AddCommand(loginCmd())
	return cmd
}
