package root

import (
	"github.com/llm-operator/cli/internal/auth"
	"github.com/llm-operator/cli/internal/version"
	"github.com/spf13/cobra"
)

// Cmd represents the base command when called without any subcommands.
var Cmd = &cobra.Command{
	Use:   "llmo",
	Short: "LLM Operator CLI",
}

// Execute adds all child commands to the root command.
func Execute() error {
	return Cmd.Execute()
}

func init() {
	Cmd.AddCommand(auth.Cmd)
	Cmd.AddCommand(version.Cmd)
	Cmd.SilenceUsage = true
}
