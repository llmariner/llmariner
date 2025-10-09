package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/llmariner/llmariner/cli/internal/admin/org"
	"github.com/llmariner/llmariner/cli/internal/admin/project"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/llmariner/llmariner/cli/internal/ui"
	uv1 "github.com/llmariner/user-manager/api/v1"
	"github.com/llmariner/user-manager/pkg/role"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/users"
)

// Cmd is the root command for users.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "users",
		Short:              "users commands",
		Aliases:            []string{"user"},
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}

	cmd.AddCommand(listCmd())

	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd.Context())
		},
	}
}

func list(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := uv1.ListUsersRequest{}
	var resp uv1.ListUsersResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	orgs, err := org.ListOrganizations(env)
	if err != nil {
		return err
	}
	orgsByID := make(map[string]*uv1.Organization)
	for _, o := range orgs {
		orgsByID[o.Id] = o
	}

	projectsByID := make(map[string]*uv1.Project)
	for _, o := range orgs {
		projs, err := project.ListProjects(env, o.Title)
		if err != nil {
			return fmt.Errorf("list projects for org %q: %s", o.Id, err)
		}
		for _, p := range projs {
			projectsByID[p.Id] = p
		}
	}

	tbl := table.New("ID", "Organization Role Bindings", "Project Role Bindings")
	ui.FormatTable(tbl)

	for _, u := range resp.Users {
		var orbs []string
		for _, orb := range u.OrganizationRoleBindings {
			r, ok := role.OrganizationRoleToString(orb.Role)
			if !ok {
				return fmt.Errorf("invalid role %q", orb.Role)
			}

			o, ok := orgsByID[orb.OrganizationId]
			if !ok {
				return fmt.Errorf("organization %q not found", orb.OrganizationId)
			}

			orbs = append(orbs, fmt.Sprintf("%s:%s", o.Title, r))
		}
		var prbs []string
		for _, prb := range u.ProjectRoleBindings {
			r, ok := role.ProjectRoleToString(prb.Role)
			if !ok {
				return fmt.Errorf("invalid role %q", prb.Role)
			}

			p, ok := projectsByID[prb.ProjectId]
			if !ok {
				return fmt.Errorf("project %q not found", prb.ProjectId)
			}

			prbs = append(prbs, fmt.Sprintf("%s:%s", p.Title, r))
		}
		tbl.AddRow(u.Id, strings.Join(orbs, ", "), strings.Join(prbs, ", "))
	}
	tbl.Print()

	return nil
}
