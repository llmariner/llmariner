package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	jv1 "github.com/llmariner/job-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/k8s"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	itime "github.com/llmariner/llmariner/cli/internal/time"
	"github.com/llmariner/llmariner/cli/internal/ui"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/batch/jobs"
)

// Cmd is the root command for batch jobs.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "jobs",
		Short:              "Batch Job commands",
		Aliases:            []string{"job"},
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(cancelCmd())
	cmd.AddCommand(logsCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var (
		envs   []string
		fpaths []string
		opts   createOpts
	)
	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.scripts = make(map[string][]byte, len(fpaths))
			for _, fpath := range fpaths {
				data, err := os.ReadFile(fpath)
				if err != nil {
					return err
				}
				opts.scripts[filepath.Base(fpath)] = data
			}

			if len(envs) > 0 {
				opts.envs = make(map[string]string, len(envs))
				for _, e := range envs {
					ss := strings.SplitN(e, "=", 2)
					if len(ss) != 2 {
						return fmt.Errorf("invalid env format: %q", e)
					}
					opts.envs[ss[0]] = ss[1]
				}
			}
			return create(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.image, "image", "", "Image for the batch job")
	cmd.Flags().StringVar(&opts.command, "command", "", "Command to be run by the batch job")
	cmd.Flags().StringArrayVar(&envs, "env", nil, "Environment variables used within the batch job (e.g., MY_ENV=somevalue)")
	cmd.Flags().StringArrayVar(&opts.fileIDs, "file-id", nil, "Data file id that will be downloaded to the job container")
	cmd.Flags().StringArrayVar(&fpaths, "from-file", nil, "Specify the path to a file that will be loaded as a job script")
	cmd.Flags().Int32Var(&opts.gpuCount, "gpu", 0, "Number of GPUs")
	cmd.Flags().Int32Var(&opts.workerCount, "workers", 0, "Number of workers for PyTorch DDP")

	_ = cmd.MarkFlagRequired("image")
	_ = cmd.MarkFlagRequired("command")
	_ = cmd.MarkFlagRequired("from-file")
	_ = cmd.MarkFlagFilename("from-file")
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
	cmd := &cobra.Command{
		Use:               "get <ID>",
		Args:              validateIDArg,
		ValidArgsFunction: compJobIDs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return get(cmd.Context(), args[0])
		},
	}
	return cmd
}

func deleteCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:               "delete <ID>",
		Args:              validateIDArg,
		ValidArgsFunction: compJobIDs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), args[0], force)
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip the confirmation prompt")
	return cmd
}

func cancelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "cancel <ID>",
		Args:              validateIDArg,
		ValidArgsFunction: compJobIDs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cancel(cmd.Context(), args[0])
		},
	}
	return cmd
}

func logsCmd() *cobra.Command {
	var follow bool
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

type createOpts struct {
	image       string
	command     string
	scripts     map[string][]byte
	fileIDs     []string
	envs        map[string]string
	gpuCount    int32
	workerCount int32
}

func create(ctx context.Context, opts createOpts) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := jv1.CreateBatchJobRequest{
		Image:     opts.image,
		Command:   opts.command,
		Scripts:   opts.scripts,
		Envs:      opts.envs,
		DataFiles: opts.fileIDs,
	}
	if opts.gpuCount > 0 {
		req.Resources = &jv1.BatchJob_Resources{
			GpuCount: opts.gpuCount,
		}
	}
	if opts.workerCount > 0 {
		req.Kind = &jv1.BatchJob_Kind{
			Kind: &jv1.BatchJob_Kind_Pytorch{
				Pytorch: &jv1.PyTorchJob{
					WorkerCount: opts.workerCount,
				},
			},
		}
	}

	var resp jv1.BatchJob
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}
	fmt.Printf("Created the batch job (ID: %q).\n", resp.Id)
	return nil
}

func list(ctx context.Context) error {
	nbs, err := listBatchJobs(ctx)
	if err != nil {
		return err
	}

	tbl := table.New("ID", "Image", "Status", "Age")
	ui.FormatTable(tbl)
	for _, j := range nbs {
		var age string
		if j.CreatedAt > 0 {
			age = itime.ToAge(time.Unix(j.CreatedAt, 0))
		}
		tbl.AddRow(
			j.Id,
			j.Image,
			j.Status,
			age,
		)
	}
	tbl.Print()
	return nil
}

func get(ctx context.Context, id string) error {
	return sendRequestAndPrintBatchJob(ctx, http.MethodGet, fmt.Sprintf("%s/%s", path, id), &jv1.GetBatchJobRequest{})
}

func delete(ctx context.Context, id string, force bool) error {
	if !force {
		p := ui.NewPrompter()
		s := &survey.Confirm{
			Message: fmt.Sprintf("Delete batch job %q?", id),
			Default: false,
		}
		var ok bool
		if err := p.Ask(s, &ok); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return sendRequestAndPrintBatchJob(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", path, id), &jv1.GetBatchJobRequest{})
}

func cancel(ctx context.Context, id string) error {
	return sendRequestAndPrintBatchJob(ctx, http.MethodPost, fmt.Sprintf("%s/%s/cancel", path, id), &jv1.CancelBatchJobRequest{})
}

func validateIDArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<ID> is required argument")
	}
	return nil
}

func compJobIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	list, err := listBatchJobs(cmd.Context())
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

func listBatchJobs(ctx context.Context) ([]*jv1.BatchJob, error) {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil, err
	}
	var nbs []*jv1.BatchJob
	var after string
	for {
		req := jv1.ListBatchJobsRequest{
			After: after,
		}
		var resp jv1.ListBatchJobsResponse
		if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
			return nil, err
		}
		nbs = append(nbs, resp.Jobs...)
		if !resp.HasMore {
			break
		}
		after = resp.Jobs[len(resp.Jobs)-1].Id
	}
	return nbs, nil
}

func logs(ctx context.Context, id string, follow bool) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	var job jv1.BatchJob
	if err := ihttp.NewClient(env).Send(http.MethodGet, fmt.Sprintf("%s/%s", path, id), &jv1.GetBatchJobRequest{}, &job); err != nil {
		return err
	}

	pods, err := k8s.ListPodsForJob(ctx, job.ClusterId, job.KubernetesNamespace, job.Id)
	if err != nil {
		return err
	}
	return k8s.StreamPodsLogs(ctx, job.ClusterId, follow, pods...)
}

func sendRequestAndPrintBatchJob(ctx context.Context, method, path string, req any) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	var resp jv1.BatchJob
	if err := ihttp.NewClient(env).Send(method, path, req, &resp); err != nil {
		return err
	}
	return printBatchJob(&resp)
}

func printBatchJob(nb *jv1.BatchJob) error {
	b, err := json.MarshalIndent(&nb, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
