package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/k8s"
	"github.com/llm-operator/cli/internal/runtime"
	"github.com/llm-operator/cli/internal/ui"
	jv1 "github.com/llm-operator/job-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1t "k8s.io/client-go/kubernetes/typed/core/v1"
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
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(cancelCmd())
	cmd.AddCommand(logsCmd())
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

func cancelCmd() *cobra.Command {
	var (
		id string
	)
	cmd := &cobra.Command{
		Use:  "cancel",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cancel(cmd.Context(), id)
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "ID of the job")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func logsCmd() *cobra.Command {
	var (
		id string
	)
	cmd := &cobra.Command{
		Use:  "logs",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logs(cmd.Context(), id)
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "ID of the job")
	_ = cmd.MarkFlagRequired("pod")
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

func logs(ctx context.Context, id string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	c := ihttp.NewClient(env)

	var job jv1.Job
	if err := c.Send(http.MethodGet, fmt.Sprintf("%s/%s", path, id), &jv1.GetJobRequest{}, &job); err != nil {
		return err
	}
	namespace := job.KubernetesNamespace

	kc, err := k8s.NewClient(env)
	if err != nil {
		return err
	}
	podClient := kc.CoreV1().Pods(namespace)
	resp, err := podClient.List(ctx, metav1.ListOptions{
		// This is an implicit assumption that the pod name is "job-<job_id>".
		LabelSelector: fmt.Sprintf("job-name=job-%s", id),
	})
	if err != nil {
		return err
	}
	if len(resp.Items) == 0 {
		return fmt.Errorf("no pod found for the job %q", id)
	}

	// Choose the latest pod or the last failed job.
	var latestPod *corev1.Pod
	var lastFailed *corev1.Pod
	for _, pod := range resp.Items {
		if latestPod == nil || pod.CreationTimestamp.After(latestPod.CreationTimestamp.Time) {
			latestPod = &pod
		}

		if pod.Status.Phase != corev1.PodFailed {
			continue
		}
		if lastFailed == nil || pod.CreationTimestamp.After(lastFailed.CreationTimestamp.Time) {
			lastFailed = &pod
		}
	}

	if lastFailed != nil {
		return podLog(ctx, podClient, lastFailed)
	}

	return podLog(ctx, podClient, latestPod)
}

func podLog(ctx context.Context, client corev1t.PodInterface, pod *corev1.Pod) error {
	req := client.GetLogs(pod.Name, &corev1.PodLogOptions{
		Follow: true,
	})
	stream, err := req.Stream(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = stream.Close()
	}()
	_, err = io.Copy(os.Stdout, stream)
	if err != nil {
		return err
	}
	return nil
}
