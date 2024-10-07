package workspaces

import (
	"github.com/llmariner/llmariner/cli/internal/workspaces/notebooks"
	"github.com/spf13/cobra"
)

// Cmd is the root command for workspace.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "workspace",
		Short:              "workspace commands",
		Aliases:            []string{"ws"},
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(notebooks.Cmd())
	return cmd
}
