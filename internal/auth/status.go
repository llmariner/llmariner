package auth

import (
	"github.com/llm-operator/cli/internal/accesstoken"
	"github.com/llm-operator/cli/internal/ui"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Login status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := ui.NewPrompter()

			if _, err := accesstoken.LoadToken(); err != nil {
				p.Warn("Not logged in.")
				return nil
			}

			p.Printf("Logged in.\n")
			p.Printf("- Token file location: %s\n", accesstoken.TokenFilePath())
			return nil
		},
	}
}
