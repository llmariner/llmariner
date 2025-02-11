package jobs

import (
	"github.com/llmariner/llmariner/cli/internal/hidden/jobs/clusters"
	"github.com/spf13/cobra"
)

// Cmd is the root command for jobs.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "jobs",
		Short:              "Jobs commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.Hidden = true
	cmd.AddCommand(clusters.Cmd())
	return cmd
}
