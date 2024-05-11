package models

import (
	"context"
	"net/http"
	"time"

	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/runtime"
	"github.com/llm-operator/cli/internal/ui"
	mv1 "github.com/llm-operator/model-manager/api/v1"
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

	var req mv1.ListModelsRequest
	var resp mv1.ListModelsResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	tbl := table.New("ID", "Owned By", "Finished At")
	ui.FormatTable(tbl)

	for _, m := range resp.Data {
		tbl.AddRow(
			m.Id,
			m.OwnedBy,
			time.Unix(m.Created, 0).Format(time.RFC3339),
		)
	}

	tbl.Print()

	return nil
}
