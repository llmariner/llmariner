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
	snapshotPath = "/clustertelemetry/clustersnapshots"
	gpuUsagePath = "/clustertelemetry/gpu-usages"
)

// Cmd is the root command for usage.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "clustert-telemetry",
		Short:              "Cluster Telemetry commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.Hidden = true
	cmd.AddCommand(listClusterSnapshotsCmd())
	cmd.AddCommand(listGPUUsagesCmd())
	return cmd
}

func listClusterSnapshotsCmd() *cobra.Command {
	var groupByStr string
	cmd := &cobra.Command{
		Use:   "list-cluster-snapshots",
		Short: "List cluster snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupBy, err := toGroupByEnum(groupByStr)
			if err != nil {
				return err
			}
			return listClusterSnapshots(cmd.Context(), groupBy)
		},
	}

	cmd.Flags().StringVar(&groupByStr, "group-by", "", "Group snapshots with a specified key (optional). Either 'cluster' or 'product'.")
	return cmd
}

func listGPUUsagesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-gpu-usages",
		Short: "List GPU usages",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listGPUUsages(cmd.Context())
		},
	}
}

func listClusterSnapshots(ctx context.Context, groupBy cv1.ListClusterSnapshotsRequest_GroupBy) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	req := cv1.ListClusterSnapshotsRequest{
		GroupBy: groupBy,
	}
	var resp cv1.ListClusterSnapshotsResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, snapshotPath, &req, &resp); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&resp, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func listGPUUsages(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	req := cv1.ListGpuUsagesRequest{}
	var resp cv1.ListGpuUsagesResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, gpuUsagePath, &req, &resp); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&resp, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func toGroupByEnum(groupByStr string) (cv1.ListClusterSnapshotsRequest_GroupBy, error) {
	switch groupByStr {
	case "":
		return cv1.ListClusterSnapshotsRequest_GROUP_BY_UNSPECIFIED, nil
	case "cluster":
		return cv1.ListClusterSnapshotsRequest_GROUP_BY_CLUSTER, nil
	case "product":
		return cv1.ListClusterSnapshotsRequest_GROUP_BY_PRODUCT, nil
	default:
		return cv1.ListClusterSnapshotsRequest_GROUP_BY_UNSPECIFIED, fmt.Errorf("invalid repository type %q. Must be 'cluster' or 'product'", groupByStr)
	}
}
