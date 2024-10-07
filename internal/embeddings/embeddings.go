package embeddings

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	ihttp "github.com/llmariner/llmariner/internal/http"
	"github.com/llmariner/llmariner/internal/runtime"
	iv1 "github.com/llmariner/inference-manager/api/v1"
	"github.com/spf13/cobra"
)

const (
	path = "/embeddings"
)

// Cmd is the root command for embeddings.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "embeddings",
		Short:              "Embeddings commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var req iv1.CreateEmbeddingRequest

	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), &req)
		},
	}
	cmd.Flags().StringVar(&req.Input, "input", "", "Iinput text to embed")
	cmd.Flags().StringVar(&req.Model, "model", "", "Model to be used")
	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("model")
	return cmd
}

func create(ctx context.Context, req *iv1.CreateEmbeddingRequest) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	var em iv1.Embeddings
	if err := ihttp.NewClient(env).Send(http.MethodPost, path, &req, &em); err != nil {
		return err
	}

	b, err := json.MarshalIndent(&em, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
