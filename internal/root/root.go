package root

import (
	"os"

	"github.com/llm-operator/cli/internal/admin"
	"github.com/llm-operator/cli/internal/auth"
	"github.com/llm-operator/cli/internal/batch"
	"github.com/llm-operator/cli/internal/chat"
	"github.com/llm-operator/cli/internal/context"
	"github.com/llm-operator/cli/internal/finetuning"
	"github.com/llm-operator/cli/internal/legacy"
	"github.com/llm-operator/cli/internal/models"
	"github.com/llm-operator/cli/internal/storage"
	"github.com/llm-operator/cli/internal/ui"
	"github.com/llm-operator/cli/internal/version"
	"github.com/llm-operator/cli/internal/workspaces"
	"github.com/spf13/cobra"
)

// Cmd represents the base command when called without any subcommands.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "llmo",
		Short:              "LLM Operator CLI",
		DisableFlagParsing: true,
	}
	cmd.PersistentFlags().StringVar(&ui.Color, "color", string(ui.ColorAuto), "Control color output. Available options are 'auto', 'always' and 'never'.")

	cmd.AddCommand(auth.Cmd())
	cmd.AddCommand(admin.Cmd())
	cmd.AddCommand(chat.Cmd())
	cmd.AddCommand(context.Cmd())
	cmd.AddCommand(finetuning.Cmd())
	cmd.AddCommand(workspaces.Cmd())
	cmd.AddCommand(batch.Cmd())
	cmd.AddCommand(models.Cmd())
	cmd.AddCommand(storage.Cmd())
	cmd.AddCommand(version.Cmd())

	if os.Getenv("LLMO_DEBUG") == "true" {
		cmd.AddCommand(legacy.Cmd())
	}

	cmd.SilenceUsage = true

	return cmd
}

// Execute adds all child commands to the root command.
func Execute() error {
	return Cmd().Execute()
}
