package auth

import (
	"github.com/llmariner/llmariner/internal/accesstoken"
	"github.com/llmariner/llmariner/internal/runtime"
	"github.com/llmariner/llmariner/internal/ui"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Login status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := ui.NewPrompter()

			if _, err := runtime.NewEnv(cmd.Context()); err != nil {
				p.Warn("Not logged in.")
				return nil
			}

			p.Printf("Logged in.\n")
			p.Printf("- Token file location: %s\n", accesstoken.TokenFilePath())
			return nil
		},
	}
}
