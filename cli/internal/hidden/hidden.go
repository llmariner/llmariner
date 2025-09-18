package hidden

import (
	"github.com/llmariner/llmariner/cli/internal/hidden/apiusage"
	"github.com/llmariner/llmariner/cli/internal/hidden/clustertelemetry"
	"github.com/llmariner/llmariner/cli/internal/hidden/inference"
	"github.com/llmariner/llmariner/cli/internal/hidden/jobs"
	"github.com/llmariner/llmariner/cli/internal/hidden/tokenize"
	"github.com/spf13/cobra"
)

// Cmd is the root command for hidden.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "hidden",
		Short:              "Hidden commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.Hidden = true
	cmd.AddCommand(jobs.Cmd())
	cmd.AddCommand(inference.Cmd())
	cmd.AddCommand(clustertelemetry.Cmd())
	cmd.AddCommand(apiusage.Cmd())
	cmd.AddCommand(tokenize.Cmd())
	return cmd
}
