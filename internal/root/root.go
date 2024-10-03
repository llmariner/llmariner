package root

import (
	"os"

	"github.com/llmariner/cli/internal/admin"
	"github.com/llmariner/cli/internal/auth"
	"github.com/llmariner/cli/internal/batch"
	"github.com/llmariner/cli/internal/chat"
	"github.com/llmariner/cli/internal/context"
	"github.com/llmariner/cli/internal/embeddings"
	"github.com/llmariner/cli/internal/finetuning"
	"github.com/llmariner/cli/internal/hidden"
	"github.com/llmariner/cli/internal/legacy"
	"github.com/llmariner/cli/internal/models"
	"github.com/llmariner/cli/internal/storage"
	"github.com/llmariner/cli/internal/ui"
	"github.com/llmariner/cli/internal/version"
	"github.com/llmariner/cli/internal/workspaces"
	"github.com/spf13/cobra"
)

// Cmd represents the base command when called without any subcommands.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "llma",
		Short:              "LLMariner CLI",
		DisableFlagParsing: true,
	}
	cmd.PersistentFlags().StringVar(&ui.Color, "color", string(ui.ColorAuto), "Control color output. Available options are 'auto', 'always' and 'never'.")

	cmd.AddCommand(auth.Cmd())
	cmd.AddCommand(admin.Cmd())
	cmd.AddCommand(chat.Cmd())
	cmd.AddCommand(context.Cmd())
	cmd.AddCommand(embeddings.Cmd())
	cmd.AddCommand(finetuning.Cmd())
	cmd.AddCommand(workspaces.Cmd())
	cmd.AddCommand(batch.Cmd())
	cmd.AddCommand(models.Cmd())
	cmd.AddCommand(storage.Cmd())
	cmd.AddCommand(version.Cmd())
	cmd.AddCommand(hidden.Cmd())

	if os.Getenv("LLMA_DEBUG") == "true" {
		cmd.AddCommand(legacy.Cmd())
	}

	cmd.SilenceUsage = true

	return cmd
}

// Execute adds all child commands to the root command.
func Execute() error {
	return Cmd().Execute()
}
