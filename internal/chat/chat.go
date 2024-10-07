package chat

import (
	"github.com/llmariner/llmariner/internal/chat/completions"
	"github.com/spf13/cobra"
)

// Cmd is the root command for chat.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "chat",
		Short:              "Chat commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(completions.Cmd())
	return cmd
}
