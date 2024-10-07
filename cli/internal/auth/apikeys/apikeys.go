package apikeys

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/llmariner/llmariner/cli/internal/admin/org"
	"github.com/llmariner/llmariner/cli/internal/admin/project"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/llmariner/llmariner/cli/internal/ui"
	uv1 "github.com/llmariner/user-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	pathPattern = "/organizations/%s/projects/%s/api_keys"
)

// Cmd is the root command for apikeys.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "api-keys",
		Short:              "API Keys commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(deleteCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var (
		orgTitle     string
		projectTitle string
	)
	cmd := &cobra.Command{
		Use:  "create <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), args[0], orgTitle, projectTitle)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Title of the organization. The organization in the current context is used if not specified.")
	cmd.Flags().StringVarP(&projectTitle, "project-title", "p", "", "Title of the project. The project in the current context is used if not specified.")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func listCmd() *cobra.Command {
	var (
		orgTitle     string
		projectTitle string
	)
	cmd := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd.Context(), orgTitle, projectTitle)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Title of the organization. The organization in the current context is used if not specified.")
	cmd.Flags().StringVarP(&projectTitle, "project-title", "p", "", "Title of the project. The project in the current context is used if not specified.")
	return cmd
}

func deleteCmd() *cobra.Command {
	var (
		orgTitle     string
		projectTitle string
	)
	cmd := &cobra.Command{
		Use:  "delete",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), args[0], orgTitle, projectTitle)
		},
	}
	cmd.Flags().StringVarP(&orgTitle, "organization-title", "o", "", "Title of the organization. The organization in the current context is used if not specified.")
	cmd.Flags().StringVarP(&projectTitle, "project-title", "p", "", "Title of the project. The project in the current context is used if not specified.")
	return cmd
}

func create(ctx context.Context, name, orgTitle, projectTitle string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	orgID, projectID, err := findOrgAndProject(env, orgTitle, projectTitle)
	if err != nil {
		return err
	}

	req := &uv1.CreateAPIKeyRequest{
		Name:           name,
		OrganizationId: orgID,
		ProjectId:      projectID,
	}
	var resp uv1.APIKey
	path := fmt.Sprintf(pathPattern, orgID, projectID)
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Created the API key. Secret: %s\n", resp.Secret)
	return nil
}

func list(ctx context.Context, orgTitle, projectTitle string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	orgID, projectID, err := findOrgAndProject(env, orgTitle, projectTitle)
	if err != nil {
		return err
	}

	req := &uv1.ListAPIKeysRequest{
		OrganizationId: orgID,
		ProjectId:      projectID,
	}
	var resp uv1.ListAPIKeysResponse
	path := fmt.Sprintf(pathPattern, orgID, projectID)
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, req, &resp); err != nil {
		return err
	}

	tbl := table.New("Name", "Owner", "Created At")
	ui.FormatTable(tbl)

	for _, k := range resp.Data {
		tbl.AddRow(k.Name, k.User.Id, time.Unix(k.CreatedAt, 0).Format(time.RFC3339))
	}

	tbl.Print()

	return nil
}

func delete(ctx context.Context, name, orgTitle, projectTitle string) error {
	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Delete API key %q?", name),
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

	orgID, projectID, err := findOrgAndProject(env, orgTitle, projectTitle)
	if err != nil {
		return err
	}

	key, found, err := findKeyByName(ctx, env, name, orgID, projectID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("API key %q not found", name)
	}

	req := &uv1.DeleteAPIKeyRequest{
		Id:             key.Id,
		OrganizationId: orgID,
		ProjectId:      projectID,
	}
	var resp uv1.DeleteAPIKeyResponse
	path := fmt.Sprintf(pathPattern, orgID, projectID)
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, key.Id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Deleted the API key (ID: %q).\n", key.Id)

	return nil
}

func findOrgAndProject(env *runtime.Env, orgTitle, projectTitle string) (string, string, error) {
	orgID, err := findOrgID(env, orgTitle)
	if err != nil {
		return "", "", err
	}
	projectID, err := findProjectID(env, orgTitle, projectTitle)
	if err != nil {
		return "", "", err
	}

	return orgID, projectID, nil
}

func findProjectID(env *runtime.Env, orgTitle, projectTitle string) (string, error) {
	if projectTitle == "" {
		pid := env.Config.Context.ProjectID
		if pid == "" {
			return "", fmt.Errorf("--project-title flag must be specified or the project must be specified by 'llma context set'")
		}
		return pid, nil
	}

	project, found, err := project.FindProjectByTitle(env, projectTitle, orgTitle)
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("project not found in organization")
	}
	return project.Id, nil
}

func findOrgID(env *runtime.Env, orgTitle string) (string, error) {
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
		return "", fmt.Errorf("organization not found")
	}
	return org.Id, nil
}

func findKeyByName(
	ctx context.Context,
	env *runtime.Env,
	name string,
	orgID string,
	projectID string,
) (*uv1.APIKey, bool, error) {
	req := &uv1.ListAPIKeysRequest{
		OrganizationId: orgID,
		ProjectId:      projectID,
	}
	var resp uv1.ListAPIKeysResponse
	path := fmt.Sprintf(pathPattern, orgID, projectID)
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return nil, false, err
	}

	for _, k := range resp.Data {
		if k.Name == name {
			return k, true, nil
		}
	}
	return nil, false, nil
}

func validateNameArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<NAME> is required argument")
	}
	return nil
}
