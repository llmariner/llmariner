package auth

import "github.com/spf13/cobra"

// Cmd is the root command for auth.
var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "auth",
}

func init() {
	Cmd.AddCommand(loginCmd())
	Cmd.SilenceUsage = true
}
