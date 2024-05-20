package finetuning

import (
	"github.com/llm-operator/cli/internal/finetuning/jobs"
	"github.com/spf13/cobra"
)

// Cmd is the root command for fine-tuning.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "fine-tuning",
		Short:              "Fine tuning commands",
		Aliases:            []string{"ft"},
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(jobs.Cmd())
	return cmd
}
