package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jv1 "github.com/llmariner/job-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/k8s"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	itime "github.com/llmariner/llmariner/cli/internal/time"
	"github.com/llmariner/llmariner/cli/internal/ui"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
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
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(cancelCmd())
	cmd.AddCommand(execCmd())
	cmd.AddCommand(logsCmd())
	return cmd
}

func createCmd() *cobra.Command {
	req := jv1.CreateJobRequest{
		Resources: &jv1.Job_Resources{},
	}
	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), &req)
		},
	}

	cmd.Flags().StringVar(&req.Model, "model", "", "Base model to fine-tune")
	cmd.Flags().StringVar(&req.TrainingFile, "training-file-id", "", "Training file ID")
	cmd.Flags().StringVar(&req.Suffix, "suffix", "", "Suffix for the fine-tuned model name")
	cmd.Flags().Int32Var(&req.Resources.GpuCount, "gpu", 1, "Number of GPUs to use")

	_ = cmd.MarkFlagRequired("model")
	_ = cmd.MarkFlagRequired("training-file-id")
	_ = cmd.MarkFlagRequired("suffix")

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
	return &cobra.Command{
		Use:               "get <ID>",
		Args:              validateIDArg,
		ValidArgsFunction: compJobIDs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(cmd.Context(), args[0])
		},
	}
}

func cancelCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "cancel <ID>",
		Args:              validateIDArg,
		ValidArgsFunction: compJobIDs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cancel(cmd.Context(), args[0])
		},
	}
}

func execCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "exec <ID>",
		Args:              validateIDArg,
		ValidArgsFunction: compJobIDs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return exec(cmd.Context(), args[0])
		},
	}
}

func logsCmd() *cobra.Command {
	var (
		follow bool
	)
	cmd := &cobra.Command{
		Use:               "logs <ID>",
		Args:              validateIDArg,
		ValidArgsFunction: compJobIDs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logs(cmd.Context(), args[0], follow)
		},
	}
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "True if the logs should be streamed")
	return cmd
}

func compJobIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	list, err := listJobs(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var cands []string
	for _, job := range list {
		if toComplete == "" || strings.HasPrefix(job.Id, toComplete) {
			cands = append(cands, job.Id)
		}
	}
	return cands, cobra.ShellCompDirectiveNoFileComp
}

func validateIDArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<ID> is required argument")
	}
	return nil
}

func create(ctx context.Context, req *jv1.CreateJobRequest) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var resp jv1.Job
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, req, &resp); err != nil {
		return err
	}

	fmt.Printf("Created the job (ID: %q).\n", resp.Id)

	return nil
}

func list(ctx context.Context) error {
	jobs, err := listJobs(ctx)
	if err != nil {
		return err
	}

	tbl := table.New("ID", "Model", "Fine-tuned Model", "Status", "Age")
	ui.FormatTable(tbl)

	for _, j := range jobs {
		tbl.AddRow(
			j.Id,
			j.Model,
			j.FineTunedModel,
			j.Status,
			itime.ToAge(time.Unix(j.CreatedAt, 0)),
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

func cancel(ctx context.Context, id string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := jv1.CancelJobRequest{
		Id: id,
	}
	var resp jv1.Job
	if err := ihttp.NewClient(env).Send(http.MethodPost, fmt.Sprintf("%s/%s/cancel", path, id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Canceled the job (ID: %q).\n", id)

	return nil
}

func logs(ctx context.Context, id string, follow bool) error {
	job, err := getJob(ctx, id)
	if err != nil {
		return err
	}

	pods, err := k8s.ListPodsForJob(ctx, job.ClusterId, job.KubernetesNamespace, job.Id)
	if err != nil {
		return err
	}

	latestPod, lastFailed := k8s.FindLatestOrLastFailedPod(pods)
	if lastFailed != nil {
		return k8s.StreamPodsLogs(ctx, job.ClusterId, follow, *lastFailed)
	}
	return k8s.StreamPodsLogs(ctx, job.ClusterId, follow, *latestPod)
}

func exec(ctx context.Context, id string) error {
	job, err := getJob(ctx, id)
	if err != nil {
		return err
	}

	pods, err := k8s.ListPodsForJob(ctx, job.ClusterId, job.KubernetesNamespace, job.Id)
	if err != nil {
		return err
	}

	// Choose the latest pod for the job.
	var latestPod *corev1.Pod
	for _, pod := range pods {
		if latestPod == nil || pod.CreationTimestamp.After(latestPod.CreationTimestamp.Time) {
			latestPod = &pod
		}
	}

	if latestPod == nil {
		return fmt.Errorf("no pod found for the job %q", id)
	}

	if latestPod.Status.Phase != corev1.PodRunning {
		return fmt.Errorf("the pod is not running")
	}

	return k8s.ExecPod(ctx, job.ClusterId, latestPod)
}

func getJob(ctx context.Context, id string) (*jv1.Job, error) {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil, err
	}

	var req jv1.GetJobRequest
	var resp jv1.Job
	if err := ihttp.NewClient(env).Send(http.MethodGet, fmt.Sprintf("%s/%s", path, id), &req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func listJobs(ctx context.Context) ([]*jv1.Job, error) {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil, err
	}

	var jobs []*jv1.Job
	var after string
	for {
		req := jv1.ListJobsRequest{
			After: after,
		}
		var resp jv1.ListJobsResponse
		if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
			return nil, err
		}
		jobs = append(jobs, resp.Data...)
		if !resp.HasMore {
			break
		}
		after = resp.Data[len(resp.Data)-1].Id
	}

	return jobs, nil
}
