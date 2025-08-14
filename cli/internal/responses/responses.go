package responses

import (
	"context"
	"fmt"
	"io"
	"net/http"

	iv1 "github.com/llmariner/inference-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/spf13/cobra"
)

const (
	path = "/responses"
)

// Cmd is the root command for Model Responses commands.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "responses",
		Short:              "Model Responses commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var req iv1.CreateModelResponseRequest
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a model response",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), &req)
		},
	}
	cmd.Flags().StringVar(&req.Model, "model", "", "Model to be used")
	cmd.Flags().StringVar(&req.Input, "input", "", "Text input")
	cmd.Flags().Float64Var(&req.Temperature, "temperature", 1.0, "Temperature")
	cmd.Flags().Float64Var(&req.TopP, "top-p", 1.0, "Top p")

	_ = cmd.MarkFlagRequired("model")
	_ = cmd.MarkFlagRequired("input")
	return cmd
}

func create(ctx context.Context, req *iv1.CreateModelResponseRequest) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	body, err := ihttp.NewClient(env).SendRequest(http.MethodPost, path, &req)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
