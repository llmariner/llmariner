package inference

import (
	"github.com/llmariner/llmariner/cli/internal/hidden/inference/status"
	"github.com/spf13/cobra"
)

// Cmd is the root command for inference.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "inference",
		Short:              "Inference commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.Hidden = true
	cmd.AddCommand(status.Cmd())
	return cmd
}
