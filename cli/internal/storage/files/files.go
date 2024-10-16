package files

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dustin/go-humanize"
	fv1 "github.com/llmariner/file-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/llmariner/llmariner/cli/internal/ui"
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
	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Delete file %q?", id),
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

func validateIDArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<ID> is required argument")
	}
	return nil
}
