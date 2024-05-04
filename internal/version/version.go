package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// We set these values on compilation through the -ldflags flag.
	gitTag       string
	gitCommitSha string
)

// Cmd represents the version command.
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:                "version",
		Short:              "CLI version",
		Args:               cobra.NoArgs,
		DisableFlagParsing: false,
		Run: func(cmd *cobra.Command, args []string) {
			if gitTag == "" && gitCommitSha == "" {
				fmt.Println("No version associated")
				return
			}
			fmt.Printf("Version: %s, Git commit: %s\n", gitTag, gitCommitSha)
		},
	}
}
