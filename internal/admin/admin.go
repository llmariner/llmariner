package admin

import (
	"github.com/llm-operator/cli/internal/admin/clusters"
	"github.com/llm-operator/cli/internal/admin/org"
	"github.com/llm-operator/cli/internal/admin/project"
	"github.com/spf13/cobra"
)

// Cmd is the root command for admin.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "admin",
		Short:              "Admin commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(clusters.Cmd())
	cmd.AddCommand(org.Cmd())
	cmd.AddCommand(project.Cmd())
	return cmd
}
