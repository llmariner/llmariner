package notebooks

import (
	"context"
	"encoding/json"
	"errors"
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
	return cmd
}

func createCmd() *cobra.Command {
	var (
		name    string
		imgType string
	)
	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), name, imgType)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Name of the Notebook")
	cmd.Flags().StringVar(&imgType, "image-type", "jupyter-lab-base", "Type of the Notebook Image")
	_ = cmd.MarkFlagRequired("name")
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

func create(ctx context.Context, name, imageType string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := jv1.CreateNotebookRequest{
		Name: name,
		Image: &jv1.CreateNotebookRequest_Image{
			Image: &jv1.CreateNotebookRequest_Image_Type{Type: imageType},
		},
	}
	var resp jv1.Notebook
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}
	return printNotebook(&resp)
}

func list(ctx context.Context) error {
	nbs, err := listNotebooks(ctx)
	if err != nil {
		return err
	}

	tbl := table.New("ID", "Name", "Image", "Status", "Age")
	ui.FormatTable(tbl)
	for _, j := range nbs {
		var age string
		if j.StartedAt > 0 {
			age = timeToAge(time.Unix(j.StartedAt, 0))
		}
		tbl.AddRow(
			j.Id,
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

func validateNameArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<NAME> is required argument")
	}
	return nil
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

// timeToAge formats a time into an human-redable age string.
func timeToAge(t time.Time) string {
	d := time.Since(t)
	if sec := int(d.Seconds()); sec < 60 {
		return fmt.Sprintf("%ds", sec)
	} else if min := int(d.Minutes()); min < 60 {
		return fmt.Sprintf("%dm", min)
	} else if d.Hours() < 6 {
		return fmt.Sprintf("%.0fh%dm", d.Hours(), min%60)
	} else if d.Hours() < 24 {
		return fmt.Sprintf("%.0fh", d.Hours())
	} else if d.Hours() < 24*7 {
		return fmt.Sprintf("%.0fd", d.Hours()/24)
	}
	return fmt.Sprintf("%.0fy", d.Hours()/(24*365))
}
