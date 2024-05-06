package apikeys

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/llm-operator/cli/internal/accesstoken"
	"github.com/llm-operator/cli/internal/config"
	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/ui"
	uv1 "github.com/llm-operator/user-manager/api/v1"
	v1 "github.com/llm-operator/user-manager/api/v1"
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
		DisableFlagParsing: false,
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
			c, err := config.LoadOrCreate()
			if err != nil {
				return fmt.Errorf("load or create config: %s", err)
			}
			t, err := accesstoken.LoadToken()
			if err != nil {
				return fmt.Errorf("load token: %s", err)
			}
			return create(cmd.Context(), c, t, name)
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
			c, err := config.LoadOrCreate()
			if err != nil {
				return fmt.Errorf("load or create config: %s", err)
			}
			t, err := accesstoken.LoadToken()
			if err != nil {
				return fmt.Errorf("load token: %s", err)
			}
			return list(cmd.Context(), c, t)
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
			c, err := config.LoadOrCreate()
			if err != nil {
				return fmt.Errorf("load or create config: %s", err)
			}
			t, err := accesstoken.LoadToken()
			if err != nil {
				return fmt.Errorf("load token: %s", err)
			}
			return delete(cmd.Context(), c, t, name)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Name of the API key")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func create(ctx context.Context, c *config.C, t *accesstoken.T, name string) error {
	req := &uv1.CreateAPIKeyRequest{
		Name: name,
	}
	var resp uv1.APIKey
	if err := ihttp.NewClient(c, t).Send(http.MethodPost, path, &req, &resp); err != nil {
		return err
	}

	log.Printf("Created a new API key. Secret: %s\n", resp.Secret)
	return nil
}

func list(ctx context.Context, c *config.C, t *accesstoken.T) error {
	var req uv1.ListAPIKeysRequest
	var resp uv1.ListAPIKeysResponse
	if err := ihttp.NewClient(c, t).Send(http.MethodGet, path, &req, &resp); err != nil {
		return err
	}

	tbl := table.New("ID", "Name", "Created At")
	ui.FormatTable(tbl)

	for _, k := range resp.Data {
		tbl.AddRow(k.Id, k.Name, time.Unix(k.CreatedAt, 0).Format(time.RFC3339))
	}

	tbl.Print()

	return nil
}

func delete(ctx context.Context, c *config.C, t *accesstoken.T, name string) error {
	key, found, err := findKeyByName(ctx, c, t, name)
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
	if err := ihttp.NewClient(c, t).Send(http.MethodDelete, fmt.Sprintf("%s/%s", path, key.Id), &req, &resp); err != nil {
		return err
	}

	log.Printf("Deleted the API key (ID: %q).\n", key.Id)

	return nil
}

func findKeyByName(ctx context.Context, c *config.C, t *accesstoken.T, name string) (*v1.APIKey, bool, error) {
	var req uv1.ListAPIKeysRequest
	var resp uv1.ListAPIKeysResponse
	if err := ihttp.NewClient(c, t).Send(http.MethodGet, path, &req, &resp); err != nil {
		return nil, false, err
	}

	for _, k := range resp.Data {
		if k.Name == name {
			return k, true, nil
		}
	}
	return nil, false, nil
}
