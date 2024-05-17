package files

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/runtime"
	"github.com/llm-operator/cli/internal/ui"
	fv1 "github.com/llm-operator/file-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/files"
)

// Cmd is the root command for files.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "files",
		Short:              "Files commands",
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
	var (
		id string
	)
	cmd := &cobra.Command{
		Use:  "delete",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), id)
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "ID of the file to delete")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func list(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var req fv1.ListFilesRequest
	var resp fv1.ListFilesResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	tbl := table.New("ID", "Filename", "Purpose", "Size", "Created At")
	ui.FormatTable(tbl)

	for _, f := range resp.Data {
		tbl.AddRow(
			f.Id,
			f.Filename,
			f.Purpose,
			humanize.IBytes(uint64(f.Bytes)),
			time.Unix(f.CreatedAt, 0).Format(time.RFC3339),
		)
	}

	tbl.Print()

	return nil
}

func delete(ctx context.Context, id string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := &fv1.DeleteFileRequest{
		Id: id,
	}
	var resp fv1.DeleteFileResponse
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Deleted the file (ID: %q).\n", id)

	return nil
}
