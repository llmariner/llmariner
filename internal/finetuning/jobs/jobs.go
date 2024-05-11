package jobs

import (
	"context"
	"net/http"
	"time"

	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/runtime"
	"github.com/llm-operator/cli/internal/ui"
	jv1 "github.com/llm-operator/job-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/fine_tuning/jobs"
)

// Cmd is the root command for jobs.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "jobs",
		Short:              "Jobs commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: false,
	}
	cmd.AddCommand(listCmd())
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

func list(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var req jv1.ListJobsRequest
	var resp jv1.ListJobsResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	tbl := table.New("ID", "Model", "Fine-tuned Model", "Status", "Created At", "Finished At")
	ui.FormatTable(tbl)

	for _, j := range resp.Data {
		tbl.AddRow(
			j.Id,
			j.Model,
			j.FineTunedModel,
			j.Status,
			time.Unix(j.CreatedAt, 0).Format(time.RFC3339),
		)
	}

	tbl.Print()

	return nil
}
