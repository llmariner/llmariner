package clusters

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	cv1 "github.com/llmariner/cluster-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/llmariner/llmariner/cli/internal/ui"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/clusters"
)

// Cmd is the root command for clusters.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "clusters",
		Short:              "Clusters commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(listCmd())
	cmd.AddCommand(registerCmd())
	cmd.AddCommand(unregisterCmd())
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

func registerCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "register <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return register(cmd.Context(), args[0])
		},
	}
}

func unregisterCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "unregister <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return unregister(cmd.Context(), args[0])
		},
	}
}

func register(ctx context.Context, name string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := &cv1.CreateClusterRequest{
		Name: name,
	}
	var resp cv1.Cluster
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Registered the cluster (ID: %q, Registration Key: %q).\n", resp.Id, resp.RegistrationKey)

	return nil
}

func list(ctx context.Context) error {
	cs, err := listClusters(ctx)
	if err != nil {
		return err
	}

	tbl := table.New("NAME")
	ui.FormatTable(tbl)

	for _, c := range cs {
		tbl.AddRow(c.Name)
	}

	tbl.Print()

	return nil
}

func unregister(ctx context.Context, name string) error {
	id, err := getClusterIDByName(ctx, name)
	if err != nil {
		return err
	}

	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Unregister Cluster %q?", name),
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

	req := &cv1.DeleteClusterRequest{
		Id: id,
	}
	var resp cv1.DeleteClusterResponse
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Unregistered the cluster (ID: %q).\n", id)

	return nil
}

func getClusterIDByName(ctx context.Context, name string) (string, error) {
	cs, err := listClusters(ctx)
	if err != nil {
		return "", nil
	}
	for _, c := range cs {
		if c.Name == name {
			return c.Id, nil
		}
	}
	return "", fmt.Errorf("cluster %q not found", name)
}

func listClusters(ctx context.Context) ([]*cv1.Cluster, error) {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil, err
	}
	var req cv1.ListClustersRequest
	var resp cv1.ListClustersResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func validateNameArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<NAME> is required argument")
	}
	return nil
}
