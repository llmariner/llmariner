package context

import (
	"context"
	"fmt"

	"github.com/llm-operator/cli/internal/auth/org"
	"github.com/llm-operator/cli/internal/auth/project"
	"github.com/llm-operator/cli/internal/runtime"
	uv1 "github.com/llm-operator/user-manager/api/v1"
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
		Use: "get",
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(cmd.Context(), args)
		},
	}
}

func setCmd() *cobra.Command {
	return &cobra.Command{
		Use: "set",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("expected <key> <value>")
			}
			return set(cmd.Context(), args[0], args[1])
		},
	}
}

func get(ctx context.Context, args []string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var keys []string
	if len(args) == 0 {
		keys = append(keys, orgKey, projectKey)
	} else {
		keys = append(keys, args[0])
	}

	for _, key := range keys {
		switch key {
		case orgKey, organizationKey:
			oid := env.Config.Context.OrganizationID
			orgs, err := org.ListOrganizations(env)
			if err != nil {
				return err
			}
			var found *uv1.Organization
			for _, o := range orgs {
				if o.Id == oid {
					found = o
					break
				}
			}
			if found == nil {
				return fmt.Errorf("org %q not found", oid)
			}
			fmt.Printf("Organization:\n  Title: %q\n  ID: %q\n", found.Title, oid)
		case projectKey:
			pid := env.Config.Context.ProjectID
			ps, err := project.ListProjects(env, "")
			if err != nil {
				return err
			}
			var found *uv1.Project
			for _, p := range ps {
				if p.Id == pid {
					found = p
					break
				}
			}

			if found == nil {
				return fmt.Errorf("project not found")
			}
			fmt.Printf("Project:\n  Title: %q\n  ID: %q\n", found.Title, pid)
		default:
			return fmt.Errorf("unknown key %s", key)
		}
	}

	return nil
}

func set(ctx context.Context, key, value string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	switch key {
	case orgKey, organizationKey:
		o, found, err := org.FindOrgByTitle(env, value)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("organization not found")
		}
		env.Config.Context.OrganizationID = o.Id
		if err := env.Config.Save(); err != nil {
			return err
		}
	case projectKey:
		p, found, err := project.FindProjectByTitle(env, value, "")
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("project not found")
		}
		env.Config.Context.ProjectID = p.Id
		if err := env.Config.Save(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown key %s", key)
	}

	return nil
}
