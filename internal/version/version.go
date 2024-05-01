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
var Cmd = &cobra.Command{
	Use:   "version",
	Short: "CLI version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if gitTag == "" && gitCommitSha == "" {
			fmt.Println("No version associated")
			return
		}

		fmt.Printf("Version: %s, Git commit: %s\n", gitTag, gitCommitSha)
	},
}
