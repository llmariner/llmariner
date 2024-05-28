package notebooks

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
	cmd.Flags().StringVar(&id, "id", "", "ID of the notebook")
	_ = cmd.MarkFlagRequired("id")
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
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var nbs []*jv1.Notebook
	var after string
	for {
		req := jv1.ListNotebooksRequest{
			After: after,
		}
		var resp jv1.ListNotebooksResponse
		if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
			return err
		}
		nbs = append(nbs, resp.Notebooks...)
		if !resp.HasMore {
			break
		}
		after = resp.Notebooks[len(resp.Notebooks)-1].Id
	}

	tbl := table.New("ID", "Name", "Image", "Status", "Created At", "Started At", "Stopped At")
	ui.FormatTable(tbl)
	for _, j := range nbs {
		tbl.AddRow(
			j.Id,
			j.Name,
			j.Image,
			j.Status,
			time.Unix(j.CreatedAt, 0).Format(time.RFC3339),
			time.Unix(j.StartedAt, 0).Format(time.RFC3339),
			time.Unix(j.StoppedAt, 0).Format(time.RFC3339),
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

	var req jv1.GetNotebookRequest
	var resp jv1.Notebook
	if err := ihttp.NewClient(env).Send(http.MethodGet, fmt.Sprintf("%s/%s", path, id), &req, &resp); err != nil {
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
