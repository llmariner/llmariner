package storage

import (
	"github.com/llm-operator/cli/internal/storage/files"
	"github.com/llm-operator/cli/internal/storage/vectorstores"
	"github.com/spf13/cobra"
)

// Cmd is the root command for storage.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "storage",
		Short:              "Storage commands",
		Aliases:            []string{"st"},
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(files.Cmd())
	cmd.AddCommand(vectorstores.Cmd())
	return cmd
}
