package chat

import (
	"github.com/llm-operator/cli/internal/chat/completions"
	"github.com/spf13/cobra"
)

// Cmd is the root command for chat.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "chat",
		Short:              "Chat commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: false,
	}
	cmd.AddCommand(completions.Cmd())
	return cmd
}
