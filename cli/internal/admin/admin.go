package admin

import (
	"github.com/llmariner/llmariner/cli/internal/admin/clusters"
	"github.com/llmariner/llmariner/cli/internal/admin/org"
	"github.com/llmariner/llmariner/cli/internal/admin/project"
	"github.com/llmariner/llmariner/cli/internal/admin/user"
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
	cmd.AddCommand(user.Cmd())
	return cmd
}
