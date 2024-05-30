package apikeys

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/runtime"
	"github.com/llm-operator/cli/internal/ui"
	uv1 "github.com/llm-operator/user-manager/api/v1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

const (
	path = "/users/api_keys"
)

// Cmd is the root command for apikeys.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "api-keys",
		Short:              "API Keys commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(deleteCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var (
		name string
	)
	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), name)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Name of the API key")
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

func deleteCmd() *cobra.Command {
	var (
		name string
	)
	cmd := &cobra.Command{
		Use:  "delete",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd.Context(), name)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Name of the API key")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func create(ctx context.Context, name string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := &uv1.CreateAPIKeyRequest{
		Name: name,
	}
	var resp uv1.APIKey
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Created a new API key. Secret: %s\n", resp.Secret)
	return nil
}

func list(ctx context.Context) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var req uv1.ListAPIKeysRequest
	var resp uv1.ListAPIKeysResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	tbl := table.New("Name", "Owner", "Created At")
	ui.FormatTable(tbl)

	for _, k := range resp.Data {
		tbl.AddRow(k.Name, k.User.Id, time.Unix(k.CreatedAt, 0).Format(time.RFC3339))
	}

	tbl.Print()

	return nil
}

func delete(ctx context.Context, name string) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	key, found, err := findKeyByName(ctx, env, name)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("API key %q not found", name)
	}

	req := &uv1.DeleteAPIKeyRequest{
		Id: key.Id,
	}
	var resp uv1.DeleteAPIKeyResponse
	if err := ihttp.NewClient(env).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, key.Id), &req, &resp); err != nil {
		return err
	}

	fmt.Printf("Deleted the API key (ID: %q).\n", key.Id)

	return nil
}

func findKeyByName(ctx context.Context, env *runtime.Env, name string) (*uv1.APIKey, bool, error) {
	var req uv1.ListAPIKeysRequest
	var resp uv1.ListAPIKeysResponse
	if err := ihttp.NewClient(env).Send(http.MethodGet, path, &req, &resp); err != nil {
		return nil, false, err
	}

	for _, k := range resp.Data {
		if k.Name == name {
			return k, true, nil
		}
	}
	return nil, false, nil
}
