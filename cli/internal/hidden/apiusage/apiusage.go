package apiusage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	av1 "github.com/llmariner/api-usage/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/spf13/cobra"
)

const (
	path = "/api-usage/model-usage-summaries"
)

// Cmd is the root command for usage.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "api-usage",
		Short:              "API Usage commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.Hidden = true
	cmd.AddCommand(listClusterSnapshotsCmd())
	return cmd
}

func listClusterSnapshotsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-model-usage-summaries",
		Short: "List model usage summaries",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listModelUsageSummaries(cmd.Context())
		},
	}
}

func listModelUsageSummaries(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	var req av1.ListModelUsageSummariesRequest
	var resp av1.ListModelUsageSummariesResponse
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
