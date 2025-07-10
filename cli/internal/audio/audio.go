package audio

import (
	"github.com/llmariner/llmariner/cli/internal/audio/transcriptions"
	"github.com/spf13/cobra"
)

// Cmd is the root command for audio.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "audio",
		Short:              "Audio commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(transcriptions.Cmd())
	return cmd
}
