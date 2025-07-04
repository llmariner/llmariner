package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AlecAivazis/survey/v2"
	cv1 "github.com/llmariner/cluster-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/llmariner/llmariner/cli/internal/ui"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	clusterPath = "/clusters"
	pathPattern = "/clusters/%s/config"
)

// Cmd is the root command for config
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "config",
		Short:              "config commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}

	cmd.AddCommand(createCmd())
	cmd.AddCommand(getCmd())
	// TODO(kenji): Support update.
	cmd.AddCommand(deleteCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var timeSlicingGPUs int
	cmd := &cobra.Command{
		Use:  "create <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), args[0], timeSlicingGPUs)
		},
	}

	cmd.Flags().IntVar(&timeSlicingGPUs, "time-slicing-gpus", 0, "GPUs to use for time slicing")
	_ = cmd.MarkFlagRequired("time-slicing-gpus")

	return cmd
}

func getCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "get <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(cmd.Context(), args[0])
		},
	}
}

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), args[0])
		},
	}
}

func create(ctx context.Context, name string, timeSlicingGPUs int) error {
	cl, err := getClusterByName(ctx, name)
	if err != nil {
		return err
	}

	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := &cv1.CreateClusterConfigRequest{
		DevicePluginConfig: &cv1.DevicePluginConfig{
			TimeSlicing: &cv1.DevicePluginConfig_TimeSlicing{
				Gpus: int32(timeSlicingGPUs),
			},
		},
	}
	var resp cv1.ClusterConfig
	path := fmt.Sprintf(pathPattern, cl.Id)
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Created the cluster config.\n")

	return nil
}

func get(ctx context.Context, name string) error {
	cl, err := getClusterByName(ctx, name)
	if err != nil {
		return err
	}

	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	var req cv1.GetClusterConfigRequest
	var config cv1.ClusterConfig
	path := fmt.Sprintf(pathPattern, cl.Id)
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &config); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&config, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func delete(ctx context.Context, name string) error {
	cl, err := getClusterByName(ctx, name)
	if err != nil {
		return err
	}

	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Delete cluster config for %q?", name),
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

	req := &cv1.DeleteClusterConfigRequest{}
	var resp emptypb.Empty
	path := fmt.Sprintf(pathPattern, cl.Id)
	if err := ihttp.NewClient(env).Send(http.MethodDelete, path, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Unregistered the cluster config.\n")

	return nil
}

func getClusterByName(ctx context.Context, name string) (*cv1.Cluster, error) {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil, err
	}
	var req cv1.ListClustersRequest
	var resp cv1.ListClustersResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, clusterPath, &req, &resp); err != nil {
		return nil, err
	}

	for _, c := range resp.Data {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, fmt.Errorf("cluster %q not found", name)
}

func validateNameArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<NAME> is required argument")
	}
	return nil
}
