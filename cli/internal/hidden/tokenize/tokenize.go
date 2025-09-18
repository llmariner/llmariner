package tokenize

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	iv1 "github.com/llmariner/inference-manager/api/v1"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/spf13/cobra"
)

const (
	path = "/tokenize"
)

// Cmd is the root command for jobs.
func Cmd() *cobra.Command {
	req := &iv1.TokenizeRequest{}
	cmd := &cobra.Command{
		Use:   "tokenize",
		Short: "Tokenize a prompt using a specified model. Only supported by the vLLM runtime.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tokenize(cmd.Context(), req)
		},
	}
	cmd.Flags().StringVar(&req.Model, "model", "", "Model to be used")
	cmd.Flags().StringVar(&req.Prompt, "prompt", "", "Prompt to be used")
	_ = cmd.MarkFlagRequired("model")
	_ = cmd.MarkFlagRequired("prompt")

	return cmd
}

func tokenize(ctx context.Context, req *iv1.TokenizeRequest) error {
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
		return fmt.Errorf("read response: %s", err)
	}

	var resp iv1.TokenizeResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return fmt.Errorf("unmarshal response: %s", err)
	}
	fmt.Printf("Tokens: %v\n", resp.Tokens)
	fmt.Print("\n")
	return nil
}
