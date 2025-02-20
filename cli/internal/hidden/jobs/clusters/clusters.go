package clusters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	jv1 "github.com/llmariner/job-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/spf13/cobra"
)

const (
	path = "/jobs/clusters"
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
	return cmd
}

func listCmd() *cobra.Command {
	var ()
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List clusters from job-manager-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listClusters(cmd.Context())
		},
	}
	return cmd
}

func listClusters(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	var req jv1.ListClustersRequest
	var resp jv1.ListClustersResponse
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
