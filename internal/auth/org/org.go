package org

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/runtime"
	"github.com/llm-operator/cli/internal/ui"
	uv1 "github.com/llm-operator/user-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/organizations"
)

// Cmd is the root command for organizations.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "organizations",
		Short:              "organizations commands",
		Aliases:            []string{"orgs", "org"},
		Args:               cobra.NoArgs,
		DisableFlagParsing: false,
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(addMemberCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var title string
	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), title)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Title of the organization")
	_ = cmd.MarkFlagRequired("title")
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

func deleteCmd() *cobra.Command {
	var title string
	cmd := &cobra.Command{
		Use:  "delete",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), title)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Title of the organization")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func addMemberCmd() *cobra.Command {
	var title, email, role string
	cmd := &cobra.Command{
		Use:  "add-member",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			r, ok := uv1.Role_value[strings.ToUpper(role)]
			if !ok || r == 0 {
				return fmt.Errorf("invalid role %q", role)
			}
			return addMember(cmd.Context(), title, email, uv1.Role(r))
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Title of the organization")
	cmd.Flags().StringVar(&email, "email", "", "Email of the user")
	cmd.Flags().StringVar(&role, "role", "", "Role of the user (owner or reader)")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("role")
	return cmd
}

func create(ctx context.Context, title string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := uv1.CreateOrganizationRequest{
		Title: title,
	}
	var resp uv1.Organization
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Created a new organization. (ID: %s)\n", resp.Id)
	return nil
}

func list(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	orgs, err := listOrganizations(env)
	if err != nil {
		return err
	}

	tbl := table.New("Title", "Created At")
	ui.FormatTable(tbl)

	for _, o := range orgs {
		tbl.AddRow(o.Title, time.Unix(o.CreatedAt, 0).Format(time.RFC3339))
	}

	tbl.Print()

	return nil
}

func delete(ctx context.Context, title string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	org, found, err := findOrgByTitle(env, title)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("organization %q not found", title)
	}

	req := uv1.DeleteOrganizationRequest{
		Id: org.Id,
	}
	var resp uv1.DeleteOrganizationResponse
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, org.Id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Deleted the organization (ID: %q).\n", org.Id)

	return nil
}

func addMember(ctx context.Context, title, userID string, role uv1.Role) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	org, found, err := findOrgByTitle(env, title)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("organization %q not found", title)
	}

	req := uv1.AddUserToOrganizationRequest{
		User: &uv1.OrganizationUser{
			OrganizationId: org.Id,
			UserId:         userID,
			Role:           role,
		},
	}
	var resp uv1.AddUserToOrganizationResponse
	p := fmt.Sprintf("%s/%s/users:addUser", path, org.Id)
	if err := ihttp.NewClient(env).Send(http.MethodPost, p, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Added the user %q to the organization %q with role %q.\n", userID, title, role.String())

	return nil
}

func findOrgByTitle(env *runtime.Env, title string) (*uv1.Organization, bool, error) {
	orgs, err := listOrganizations(env)
	if err != nil {
		return nil, false, err
	}
	for _, o := range orgs {
		if o.Title == title {
			return o, true, nil
		}
	}
	return nil, false, nil
}

func listOrganizations(env *runtime.Env) ([]*uv1.Organization, error) {
	var req uv1.ListOrganizationsRequest
	var resp uv1.ListOrganizationsResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return resp.Organizations, nil
}
