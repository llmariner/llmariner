package notebooks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cli/browser"
	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/nbtoken"
	"github.com/llm-operator/cli/internal/runtime"
	itime "github.com/llm-operator/cli/internal/time"
	"github.com/llm-operator/cli/internal/ui"
	jv1 "github.com/llm-operator/job-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/workspaces/notebooks"
)

// Cmd is the root command for notebooks.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "notebooks",
		Short:              "Notebook commands",
		Aliases:            []string{"nbs", "nb"},
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(stopCmd())
	cmd.AddCommand(startCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(openCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var (
		envs []string
		opts createOpts
	)
	cmd := &cobra.Command{
		Use:  "create",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
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
			return create(cmd.Context(), args[0], opts)
		},
	}
	cmd.Flags().StringVar(&opts.imageType, "image-type", "jupyter-lab-base", "Type of the Notebook Image")
	cmd.Flags().StringArrayVar(&envs, "env", nil, "Environment variables used within the Notebook (e.g., MY_ENV=somevalue)")
	cmd.Flags().Int32Var(&opts.gpuCount, "gpu", 0, "Number of GPUs")
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
		Use:  "get <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			nbID, err := getNotebookIDByName(ctx, args[0])
			if err != nil {
				return err
			}
			return get(ctx, nbID)
		},
	}
	return cmd
}

func stopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "stop <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			nbID, err := getNotebookIDByName(ctx, args[0])
			if err != nil {
				return err
			}
			return stop(ctx, nbID)
		},
	}
	return cmd
}

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "start <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			nbID, err := getNotebookIDByName(ctx, args[0])
			if err != nil {
				return err
			}
			return start(ctx, nbID)
		},
	}
	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), args[0])
		},
	}
	return cmd
}

func openCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "open <NAME>",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			nbID, err := getNotebookIDByName(ctx, args[0])
			if err != nil {
				return err
			}
			return open(ctx, nbID)
		},
	}
	return cmd
}

type createOpts struct {
	imageType string
	envs      map[string]string
	gpuCount  int32
}

func create(ctx context.Context, name string, opts createOpts) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := jv1.CreateNotebookRequest{
		Name: name,
		Image: &jv1.CreateNotebookRequest_Image{
			Image: &jv1.CreateNotebookRequest_Image_Type{Type: opts.imageType},
		},
		Envs: opts.envs,
	}
	if opts.gpuCount > 0 {
		req.Resources = &jv1.Resources{
			GpuCount: opts.gpuCount,
		}
	}

	var resp jv1.Notebook
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}
	fmt.Printf("created the notebook (ID: %q).\n", resp.Id)

	return nbtoken.SaveToken(resp.Id, resp.Token)
}

func list(ctx context.Context) error {
	nbs, err := listNotebooks(ctx)
	if err != nil {
		return err
	}

	tbl := table.New("Name", "Image", "Status", "Age")
	ui.FormatTable(tbl)
	for _, j := range nbs {
		var age string
		if j.StartedAt > 0 {
			age = itime.ToAge(time.Unix(j.StartedAt, 0))
		}
		tbl.AddRow(
			j.Name,
			j.Image,
			j.Status,
			age,
		)
	}
	tbl.Print()
	return nil
}

func get(ctx context.Context, id string) error {
	return sendRequestAndPrintNotebook(ctx, http.MethodGet, fmt.Sprintf("%s/%s", path, id), &jv1.GetNotebookRequest{})
}

func stop(ctx context.Context, id string) error {
	return sendRequestAndPrintNotebook(ctx, http.MethodPost, fmt.Sprintf("%s/%s/actions:stop", path, id), &jv1.StopNotebookRequest{})
}

func start(ctx context.Context, id string) error {
	return sendRequestAndPrintNotebook(ctx, http.MethodPost, fmt.Sprintf("%s/%s/actions:start", path, id), &jv1.StartNotebookRequest{})
}

func delete(ctx context.Context, name string) error {
	id, err := getNotebookIDByName(ctx, name)
	if err != nil {
		return err
	}

	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Delete notebook %q?", name),
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
	var resp jv1.DeleteNotebookResponse
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, id), &jv1.DeleteNotebookRequest{}, &resp); err != nil {
		return err
	}
	fmt.Printf("Deleted the notebook (ID: %q).\n", id)
	return nil
}

func validateNameArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<NAME> is required argument")
	}
	return nil
}

func open(ctx context.Context, id string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	token, err := nbtoken.LoadToken(id)
	if err != nil {
		// TODO(aya): implement get token API?
		return err
	}

	var resp jv1.Notebook
	if err := ihttp.NewClient(env).Send(http.MethodGet, fmt.Sprintf("%s/%s", path, id), &jv1.GetJobRequest{}, &resp); err != nil {
		return err
	}
	if resp.Status != "running" {
		return fmt.Errorf("notebook %q is not running (status: %s)", resp.Name, resp.Status)
	}

	fmt.Println("Opening browser...")
	nbURL := fmt.Sprintf("%s/services/notebooks/%s?token=%s", env.Config.EndpointURL, id, token)
	return browser.OpenURL(nbURL)
}

func getNotebookIDByName(ctx context.Context, name string) (string, error) {
	nbs, err := listNotebooks(ctx)
	if err != nil {
		return "", nil
	}
	for _, nb := range nbs {
		if nb.Name == name {
			return nb.Id, nil
		}
	}
	return "", fmt.Errorf("notebook %q not found", name)
}

func listNotebooks(ctx context.Context) ([]*jv1.Notebook, error) {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil, err
	}
	var nbs []*jv1.Notebook
	var after string
	for {
		req := jv1.ListNotebooksRequest{
			After: after,
		}
		var resp jv1.ListNotebooksResponse
		if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
			return nil, err
		}
		nbs = append(nbs, resp.Notebooks...)
		if !resp.HasMore {
			break
		}
		after = resp.Notebooks[len(resp.Notebooks)-1].Id
	}
	return nbs, nil
}

func sendRequestAndPrintNotebook(ctx context.Context, method, path string, req any) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	var resp jv1.Notebook
	if err := ihttp.NewClient(env).Send(method, path, req, &resp); err != nil {
		return err
	}
	return printNotebook(&resp)
}

func printNotebook(nb *jv1.Notebook) error {
	b, err := json.MarshalIndent(&nb, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
