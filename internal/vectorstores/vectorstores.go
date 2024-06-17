package vectorstores

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AlecAivazis/survey/v2"
	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/runtime"
	"github.com/llm-operator/cli/internal/ui"
	vsv1 "github.com/llm-operator/vector-store-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/vector_stores"
)

// Cmd is the root command for vector stores.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "vector-stores",
		Short:              "Vector stores commands",
		Aliases:            []string{"vs"},
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	// TODO(kenji): Support create, update, and get.
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
		Use:  "delete",
		Args: validateNameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), args[0])
		},
	}
}

func list(ctx context.Context) error {
	vss, err := listVectorStores(ctx)
	if err != nil {
		return err
	}

	tbl := table.New("ID", "Name", "Status", "Created At")
	ui.FormatTable(tbl)

	for _, v := range vss {
		tbl.AddRow(
			v.Id,
			v.Name,
			v.Status,
			time.Unix(v.CreatedAt, 0).Format(time.RFC3339),
		)
	}

	tbl.Print()

	return nil
}

func delete(ctx context.Context, name string) error {
	vs, err := getVectorStoreByName(ctx, name)
	if err != nil {
		return err
	}

	p := ui.NewPrompter()
	s := &survey.Confirm{
		Message: fmt.Sprintf("Delete vector store %q?", name),
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

	req := &vsv1.DeleteVectorStoreRequest{
		Id: vs.Id,
	}
	var resp vsv1.DeleteVectorStoreResponse
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, vs.Id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Deleted the vector store (ID: %q).\n", vs.Id)

	return nil
}

func getVectorStoreByName(ctx context.Context, name string) (*vsv1.VectorStore, error) {
	vss, err := listVectorStores(ctx)
	if err != nil {
		return nil, err
	}

	for _, vs := range vss {
		if vs.Name == name {
			return vs, nil
		}
	}

	return nil, fmt.Errorf("vector store %q not found", name)
}

func listVectorStores(ctx context.Context) ([]*vsv1.VectorStore, error) {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil, err
	}

	var vss []*vsv1.VectorStore
	var after string
	for {
		req := vsv1.ListVectorStoresRequest{
			After: after,
		}
		var resp vsv1.ListVectorStoresResponse
		if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
			return nil, err
		}
		vss = append(vss, resp.Data...)
		if !resp.HasMore {
			break
		}
		after = resp.Data[len(resp.Data)-1].Id
	}

	return vss, nil
}

func validateNameArg(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("<Name> is required argument")
	}
	return nil
}
