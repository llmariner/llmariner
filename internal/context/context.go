package context

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/llmariner/llmariner/internal/admin/org"
	"github.com/llmariner/llmariner/internal/admin/project"
	"github.com/llmariner/llmariner/internal/runtime"
	"github.com/llmariner/llmariner/internal/ui"
	uv1 "github.com/llmariner/user-manager/api/v1"
	"github.com/spf13/cobra"
)

const (
	organizationKey = "organization"
	orgKey          = "org"
	projectKey      = "project"
)

// Cmd returns a new config command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "context",
		Short:              "Context commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(getCmd())
	cmd.AddCommand(setCmd())
	return cmd
}

func getCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(cmd.Context())
		},
	}
}

func setCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "set",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Set(cmd.Context())
		},
	}
}

func get(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	oid := env.Config.Context.OrganizationID
	if oid == "" {
		fmt.Printf("Organization:\n  (not selected)\n")
		return nil
	}

	orgs, err := org.ListOrganizations(env)
	if err != nil {
		return err
	}
	foundO, err := findOrgByID(orgs, oid)
	if err != nil {
		return err
	}
	fmt.Printf("Organization:\n  Title: %q\n  ID: %q\n", foundO.Title, oid)

	pid := env.Config.Context.ProjectID
	if pid == "" {
		fmt.Printf("Project:\n  (not selected)\n")
		return nil
	}

	ps, err := project.ListProjects(env, foundO.Title)
	if err != nil {
		return err
	}
	foundP, err := findProjectByID(ps, pid)
	if err != nil {
		return err
	}
	fmt.Printf("Project:\n  Title: %q\n  ID: %q\n", foundP.Title, pid)

	return nil
}

// Set sets the organization and project context.
func Set(ctx context.Context) error {
	p := ui.NewPrompter()

	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	orgs, err := org.ListOrganizations(env)
	if err != nil {
		return err
	}
	if len(orgs) == 0 {
		return fmt.Errorf("no organizations found")
	}
	defaultTitle := orgs[0].Title
	if id := env.Config.Context.OrganizationID; id != "" {
		o, err := findOrgByID(orgs, id)
		if err != nil {
			return err
		}
		defaultTitle = o.Title
	}

	var orgTitles []string
	for _, o := range orgs {
		orgTitles = append(orgTitles, o.Title)
	}

	var selectedTitle string
	if err := p.Ask(&survey.Select{
		Message: "Organization",
		Options: orgTitles,
		Default: defaultTitle,
	}, &selectedTitle); err != nil {
		return fmt.Errorf("failed to select org: %s", err)
	}

	selectedO, err := findOrgByTitle(orgs, selectedTitle)
	if err != nil {
		return err
	}

	orgChanged := env.Config.Context.OrganizationID != selectedO.Id
	env.Config.Context.OrganizationID = selectedO.Id

	projects, err := project.ListProjects(env, selectedO.Title)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		fmt.Println("No projects found. Only setting organization.")
		env.Config.Context.ProjectID = ""
		if err := env.Config.Save(); err != nil {
			return err
		}

		return nil
	}

	defaultTitle = projects[0].Title
	if id := env.Config.Context.ProjectID; id != "" && !orgChanged {
		p, err := findProjectByID(projects, id)
		if err != nil {
			return err
		}
		defaultTitle = p.Title
	}

	var projectTitles []string
	for _, p := range projects {
		projectTitles = append(projectTitles, p.Title)
	}

	if err := p.Ask(&survey.Select{
		Message: "Project",
		Options: projectTitles,
		Default: defaultTitle,
	}, &selectedTitle); err != nil {
		return fmt.Errorf("failed to select project: %s", err)
	}

	selectedP, err := findProjectByTitle(projects, selectedTitle)
	if err != nil {
		return err
	}
	env.Config.Context.ProjectID = selectedP.Id

	if err := env.Config.Save(); err != nil {
		return err
	}

	return nil
}

func findOrgByID(orgs []*uv1.Organization, id string) (*uv1.Organization, error) {
	for _, o := range orgs {
		if o.Id == id {
			return o, nil
		}
	}
	return nil, fmt.Errorf("org of ID %q not found", id)
}
func findOrgByTitle(orgs []*uv1.Organization, title string) (*uv1.Organization, error) {
	for _, o := range orgs {
		if o.Title == title {
			return o, nil
		}
	}
	return nil, fmt.Errorf("org of title %q not found", title)
}

func findProjectByID(projects []*uv1.Project, id string) (*uv1.Project, error) {
	for _, p := range projects {
		if p.Id == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("project of ID %q not found", id)
}

func findProjectByTitle(projects []*uv1.Project, title string) (*uv1.Project, error) {
	for _, p := range projects {
		if p.Title == title {
			return p, nil
		}
	}
	return nil, fmt.Errorf("project of title %q not found", title)
}
