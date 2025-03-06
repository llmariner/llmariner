package files

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	cmd.AddCommand(createLinkCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(deleteCmd())
	return cmd
}

func createLinkCmd() *cobra.Command {
	var objectPath, purpose string
	cmd := &cobra.Command{
		Use:   "create-link",
		Short: "Create a file from an object path.",
		Long:  "Create a file from an object path. A new file object will be created without upload.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.HasPrefix(objectPath, "s3://") {
				return fmt.Errorf("--object-path must not include a bucket name. For example, if a target object is located at s3://<bucket>/<path>, pass --object-path=<path>")
			}

			return createLink(cmd.Context(), objectPath, purpose)
		},
	}
	cmd.Flags().StringVar(&objectPath, "object-path", "", "Path to the object in the object storage. This does not include a bucket name.")
	cmd.Flags().StringVar(&purpose, "purpose", "", "Purpose. Either 'fine-tune' or 'assistants'.")
	_ = cmd.MarkFlagRequired("object-path")
	_ = cmd.MarkFlagRequired("purpose")
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

func createLink(ctx context.Context, objectPath, purpose string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := fv1.CreateFileFromObjectPathRequest{
		ObjectPath: objectPath,
		Purpose:    purpose,
	}
	var resp fv1.File
	if err := ihttp.NewClient(env).Send(http.MethodPost, path+":createFromObjectPath", &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Created the file (ID: %q).\n", resp.Id)
	return nil

}

func list(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var files []*fv1.File
	var after string
	for {
		req := fv1.ListFilesRequest{
			After: after,
			Limit: 1,
		}
		var resp fv1.ListFilesResponse
		if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
			return err
		}
		files = append(files, resp.Data...)
		if !resp.HasMore {
			break
		}
		after = resp.Data[len(resp.Data)-1].Id
	}

	tbl := table.New("ID", "Filename", "Purpose", "Size", "Created At")
	ui.FormatTable(tbl)

	for _, f := range files {
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
