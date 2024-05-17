package jobs

import (
	"context"
	"encoding/json"
	"fmt"
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
	cmd.AddCommand(getCmd())
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

func getCmd() *cobra.Command {
	var (
		id string
	)
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(cmd.Context(), id)
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "ID of the job")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func list(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var jobs []*jv1.Job
	var after string
	for {
		req := jv1.ListJobsRequest{
			After: after,
		}
		var resp jv1.ListJobsResponse
		if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
			return err
		}
		jobs = append(jobs, resp.Data...)
		if !resp.HasMore {
			break
		}
		after = resp.Data[len(resp.Data)-1].Id
	}

	tbl := table.New("ID", "Model", "Fine-tuned Model", "Status", "Created At", "Finished At")
	ui.FormatTable(tbl)

	for _, j := range jobs {
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

func get(ctx context.Context, id string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var req jv1.GetJobRequest
	var resp jv1.Job
	if err := ihttp.NewClient(env).Send(http.MethodGet, fmt.Sprintf("%s/%s", path, id), &req, &resp); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&resp, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
