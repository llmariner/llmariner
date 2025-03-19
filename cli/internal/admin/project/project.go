package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/llmariner/llmariner/cli/internal/admin/org"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/llmariner/llmariner/cli/internal/ui"
	uv1 "github.com/llmariner/user-manager/api/v1"
	"github.com/llmariner/user-manager/pkg/role"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	pathPattern = "/organizations/%s/projects"
)

// Cmd is the root command for projectanizations.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "projects",
		Short:              "projects commands",
		Aliases:            []string{"project"},
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}

	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(addMemberCmd())
	cmd.AddCommand(listMembersCmd())
	cmd.AddCommand(removeMemberCmd())

	return cmd
}

func createCmd() *cobra.Command {
	var orgTitle, namespace string
	cmd := &cobra.Command{
		Use:  "create <TITLE>",
		Args: validateTitleArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), args[0], orgTitle, namespace)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Organization title of the project. The organization in the current context is used if not specified.")
	cmd.Flags().StringVarP(&namespace, "kubernetes-namespace", "n", "", "Kubernetes namesapce of the project")
	_ = cmd.MarkFlagRequired("kubernetes-namespace")
	return cmd
}

func listCmd() *cobra.Command {
	var orgTitle string
	cmd := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd.Context(), orgTitle)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Organization title of the project. The organization in the current context is used if not specified.")
	return cmd
}

func getCmd() *cobra.Command {
	var orgTitle string
	cmd := &cobra.Command{
		Use:  "get <TITLE>",
		Args: validateTitleArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(cmd.Context(), args[0], orgTitle)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Organization title of the project. The organization in the current context is used if not specified.")
	return cmd
}

func deleteCmd() *cobra.Command {
	var orgTitle string
	cmd := &cobra.Command{
		Use:  "delete <TITLE>",
		Args: validateTitleArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), args[0], orgTitle)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Organization title of the project. The organization in the current context is used if not specified.")
	return cmd
}

func addMemberCmd() *cobra.Command {
	var orgTitle, email, roleStr string
	cmd := &cobra.Command{
		Use:  "add-member <TITLE>",
		Args: validateTitleArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			r, ok := role.ProjectRoleToProtoEnum(roleStr)
			if !ok {
				return fmt.Errorf("invalid role %q", roleStr)
			}
			return addMember(cmd.Context(), args[0], orgTitle, email, r)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Organization title of the project. The organization in the current context is used if not specified.")
	cmd.Flags().StringVar(&email, "email", "", "Email of the user")
	cmd.Flags().StringVar(&roleStr, "role", "", "Role of the user (owner or member)")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("role")
	return cmd
}

func listMembersCmd() *cobra.Command {
	var orgTitle string
	cmd := &cobra.Command{
		Use:  "list-members <TITLE>",
		Args: validateTitleArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listMembers(cmd.Context(), args[0], orgTitle)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Organization title of the project. The organization in the current context is used if not specified.")

	return cmd
}

func removeMemberCmd() *cobra.Command {
	var orgTitle, email string
	cmd := &cobra.Command{
		Use:  "remove-member <TITLE>",
		Args: validateTitleArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeMember(cmd.Context(), args[0], orgTitle, email)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Organization title of the project. The organization in the current context is used if not specified.")
	cmd.Flags().StringVar(&email, "email", "", "Email of the user")
	_ = cmd.MarkFlagRequired("email")
	return cmd
}

func create(ctx context.Context, title, orgTitle, namespace string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	path, orgID, err := buildPath(env, orgTitle)
	if err != nil {
		return err
	}
	req := uv1.CreateProjectRequest{
		Title:               title,
		OrganizationId:      orgID,
		KubernetesNamespace: namespace,
	}
	var resp uv1.Project
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Created the project (ID: %q).\n", resp.Id)
	return nil
}

func list(ctx context.Context, orgTitle string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	projects, err := ListProjects(env, orgTitle)
	if err != nil {
		return err
	}

	tbl := table.New("Title", "ID", "Created At")
	ui.FormatTable(tbl)

	for _, o := range projects {
		tbl.AddRow(o.Title, o.Id, time.Unix(o.CreatedAt, 0).Format(time.RFC3339))
	}

	tbl.Print()

	return nil
}

func get(ctx context.Context, title, orgTitle string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	project, found, err := FindProjectByTitle(env, title, orgTitle)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("project %q not found", title)
	}

	b, err := json.MarshalIndent(&project, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func delete(ctx context.Context, title, orgTitle string) error {
	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Delete project %q?", title),
		Default: false,
	}
	var ok bool
	if err := p.Ask(s, &ok); err != nil {
		return err
	} else if !ok {
		return nil
	}

	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	project, found, err := FindProjectByTitle(env, title, orgTitle)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("project %q not found", title)
	}

	path, orgID, err := buildPath(env, orgTitle)
	if err != nil {
		return err
	}
	req := uv1.DeleteProjectRequest{
		Id:             project.Id,
		OrganizationId: orgID,
	}
	var resp uv1.DeleteProjectResponse
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, project.Id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Deleted the project (ID: %q).\n", project.Id)

	return nil
}

func addMember(ctx context.Context, title, orgTitle, userID string, role uv1.ProjectRole) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	project, found, err := FindProjectByTitle(env, title, orgTitle)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("project %q not found", title)
	}

	path, orgID, err := buildPath(env, orgTitle)
	if err != nil {
		return err
	}
	req := uv1.CreateProjectUserRequest{
		ProjectId:      project.Id,
		OrganizationId: orgID,
		UserId:         userID,
		Role:           role,
	}
	var resp uv1.ProjectUser
	p := fmt.Sprintf("%s/%s/users", path, project.Id)
	if err := ihttp.NewClient(env).Send(http.MethodPost, p, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Added the user %q to the project %q with role %q.\n", userID, title, role.String())

	return nil
}

func listMembers(ctx context.Context, title, orgTitle string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	project, found, err := FindProjectByTitle(env, title, orgTitle)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("project %q not found", title)
	}

	path, orgID, err := buildPath(env, orgTitle)
	if err != nil {
		return err
	}
	req := uv1.ListProjectUsersRequest{
		ProjectId:      project.Id,
		OrganizationId: orgID,
	}
	var resp uv1.ListProjectUsersResponse
	p := fmt.Sprintf("%s/%s/users", path, project.Id)
	if err := ihttp.NewClient(env).Send(http.MethodGet, p, &req, &resp); err != nil {
		return err
	}

	tbl := table.New("User ID", "Role")
	ui.FormatTable(tbl)
	for _, u := range resp.Users {
		tbl.AddRow(u.UserId, u.Role.String())
	}

	tbl.Print()

	return nil
}

func removeMember(ctx context.Context, title, orgTitle, userID string) error {
	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Remove %q from project %q?", userID, title),
		Default: false,
	}
	var ok bool
	if err := p.Ask(s, &ok); err != nil {
		return err
	} else if !ok {
		return nil
	}

	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	project, found, err := FindProjectByTitle(env, title, orgTitle)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("project %q not found", title)
	}

	path, orgID, err := buildPath(env, orgTitle)
	if err != nil {
		return err
	}
	req := uv1.DeleteProjectUserRequest{
		ProjectId:      project.Id,
		OrganizationId: orgID,
		UserId:         userID,
	}
	var resp uv1.ProjectUser
	rp := fmt.Sprintf("%s/%s/users/%s", path, project.Id, userID)
	if err := ihttp.NewClient(env).Send(http.MethodDelete, rp, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Removed the user %q from the project %q.\n", userID, title)

	return nil
}

// FindProjectByTitle finds a project by title.
func FindProjectByTitle(env *runtime.Env, title, orgTitle string) (*uv1.Project, bool, error) {
	projects, err := ListProjects(env, orgTitle)
	if err != nil {
		return nil, false, err
	}
	for _, p := range projects {
		if p.Title == title {
			return p, true, nil
		}
	}
	return nil, false, nil
}

// ListProjects lists projects.
func ListProjects(env *runtime.Env, orgTitle string) ([]*uv1.Project, error) {
	path, orgID, err := buildPath(env, orgTitle)
	if err != nil {
		return nil, err
	}

	req := uv1.ListProjectsRequest{
		OrganizationId: orgID,
		IncludeSummary: true,
	}
	var resp uv1.ListProjectsResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return resp.Projects, nil
}

func buildPath(env *runtime.Env, orgTitle string) (string, string, error) {
	orgID, err := getOrgID(env, orgTitle)
	if err != nil {
		return "", "", err
	}
	return fmt.Sprintf(pathPattern, orgID), orgID, nil
}

func getOrgID(env *runtime.Env, orgTitle string) (string, error) {
	if orgTitle == "" {
		oid := env.Config.Context.OrganizationID
		if oid == "" {
			return "", fmt.Errorf("--organization-title flag must be specified or the organization must be specified by 'llma context set'")
		}
		return oid, nil
	}

	org, found, err := org.FindOrgByTitle(env, orgTitle)
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("organization %q not found", orgTitle)
	}
	return org.Id, nil
}

func validateTitleArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<TITLE> is required argument")
	}
	return nil
}
