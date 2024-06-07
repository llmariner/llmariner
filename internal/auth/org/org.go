package org

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AlecAivazis/survey/v2"
	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/runtime"
	"github.com/llm-operator/cli/internal/ui"
	uv1 "github.com/llm-operator/user-manager/api/v1"
	"github.com/llm-operator/user-manager/pkg/role"
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
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(addMemberCmd())
	cmd.AddCommand(listMembersCmd())
	cmd.AddCommand(removeMemberCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var title, namespace string
	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), title, namespace)
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
	var title, email, roleStr string
	cmd := &cobra.Command{
		Use:  "add-member",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			r, ok := role.OrganizationRoleToProtoEnum(roleStr)
			if !ok {
				return fmt.Errorf("invalid role %q. Must be 'admin' or 'reader'", roleStr)
			}
			return addMember(cmd.Context(), title, email, r)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Title of the organization")
	cmd.Flags().StringVar(&email, "email", "", "Email of the user")
	cmd.Flags().StringVar(&roleStr, "role", "", "Role of the user (owner or reader)")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("role")
	return cmd
}

func listMembersCmd() *cobra.Command {
	var title string
	cmd := &cobra.Command{
		Use:  "list-members",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listMembers(cmd.Context(), title)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Title of the organization")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("role")
	return cmd
}

func removeMemberCmd() *cobra.Command {
	var title, email string
	cmd := &cobra.Command{
		Use:  "remove-member",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeMember(cmd.Context(), title, email)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Title of the organization")
	cmd.Flags().StringVar(&email, "email", "", "Email of the user")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("email")
	return cmd
}

func create(ctx context.Context, title, namespace string) error {
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

	orgs, err := ListOrganizations(env)
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
	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Delete organization %q?", title),
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

	org, found, err := FindOrgByTitle(env, title)
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

func addMember(ctx context.Context, title, userID string, role uv1.OrganizationRole) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	org, found, err := FindOrgByTitle(env, title)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("organization %q not found", title)
	}

	req := uv1.CreateOrganizationUserRequest{
		OrganizationId: org.Id,
		UserId:         userID,
		Role:           role,
	}
	var resp uv1.OrganizationUser
	p := fmt.Sprintf("%s/%s/users", path, org.Id)
	if err := ihttp.NewClient(env).Send(http.MethodPost, p, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Added the user %q to the organization %q with role %q.\n", userID, title, role.String())

	return nil
}

func listMembers(ctx context.Context, title string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	org, found, err := FindOrgByTitle(env, title)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("organization %q not found", title)
	}

	req := uv1.ListOrganizationUsersRequest{
		OrganizationId: org.Id,
	}
	var resp uv1.ListOrganizationUsersResponse
	p := fmt.Sprintf("%s/%s/users", path, org.Id)
	if err := ihttp.NewClient(env).Send(http.MethodGet, p, &req, &resp); err != nil {
		return err
	}

	tbl := table.New("User ID", "Role")
	ui.FormatTable(tbl)
	for _, u := range resp.Users {
		r, ok := role.OrganizationRoleToString(u.Role)
		if !ok {
			return fmt.Errorf("invalid role %q", u.Role)
		}
		tbl.AddRow(u.UserId, r)
	}

	tbl.Print()

	return nil
}

func removeMember(ctx context.Context, title, userID string) error {
	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Remove %q from organization %q?", userID, title),
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

	org, found, err := FindOrgByTitle(env, title)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("organization %q not found", title)
	}

	req := uv1.DeleteOrganizationUserRequest{
		OrganizationId: org.Id,
		UserId:         userID,
	}
	var resp uv1.OrganizationUser
	rp := fmt.Sprintf("%s/%s/users/%s", path, org.Id, userID)
	if err := ihttp.NewClient(env).Send(http.MethodDelete, rp, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Removed the user %q from the organization %q.\n", userID, title)

	return nil
}

// FindOrgByTitle finds an organization by title.
func FindOrgByTitle(env *runtime.Env, title string) (*uv1.Organization, bool, error) {
	orgs, err := ListOrganizations(env)
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

// ListOrganizations lists organizations.
func ListOrganizations(env *runtime.Env) ([]*uv1.Organization, error) {
	var req uv1.ListOrganizationsRequest
	var resp uv1.ListOrganizationsResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return resp.Organizations, nil
}
