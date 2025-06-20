package clustertelemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	cv1 "github.com/llmariner/cluster-monitor/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/spf13/cobra"
)

const (
	path = "/clustertelemetry/clustersnapshots"
)

// Cmd is the root command for usage.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "clustertelemetry",
		Short:              "Cluster Telemetry commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.Hidden = true
	cmd.AddCommand(listClusterSnapshotsCmd())
	return cmd
}

func listClusterSnapshotsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-cluster-snapshots",
		Short: "List cluster snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listClusterSnapshots(cmd.Context())
		},
	}
}

func listClusterSnapshots(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	var req cv1.ListClusterSnapshotsRequest
	var resp cv1.ListClusterSnapshotsResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&resp, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
