package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/AlecAivazis/survey/v2"
	ihttp "github.com/llmariner/llmariner/internal/http"
	"github.com/llmariner/llmariner/internal/runtime"
	"github.com/llmariner/llmariner/internal/ui"
	mv1 "github.com/llmariner/model-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/models"
)

// Cmd is the root command for models.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "models",
		Short:              "Models commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(listCmd())
	cmd.AddCommand(deleteCmd())
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

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "delete <ID>",
		Args: validateIDArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), args[0])
		},
	}
}

func list(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var req mv1.ListModelsRequest
	var resp mv1.ListModelsResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	// Sort models by names.
	var ms []*mv1.Model
	ms = append(ms, resp.Data...)
	sort.Slice(ms, func(i, j int) bool {
		return ms[i].Id < ms[j].Id
	})

	tbl := table.New("ID", "Owned By", "Created At")
	ui.FormatTable(tbl)

	for _, m := range ms {
		tbl.AddRow(
			m.Id,
			m.OwnedBy,
			time.Unix(m.Created, 0).Format(time.RFC3339),
		)
	}

	tbl.Print()

	return nil
}

func delete(ctx context.Context, id string) error {
	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Delete model %q?", id),
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

	req := &mv1.DeleteModelRequest{
		Id: id,
	}
	var resp mv1.DeleteModelResponse
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Deleted the model (ID: %q).\n", id)

	return nil
}

func validateIDArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<ID> is required argument")
	}
	return nil
}
