package completions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/runtime"
	iv1 "github.com/llm-operator/inference-manager/api/v1"
	"github.com/llm-operator/inference-manager/common/pkg/sse"
	"github.com/spf13/cobra"
)

const (
	path = "/completions"
)

// Cmd is the root command for completions.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "completions",
		Short:              "Completions commands",
		Args:               cobra.NoArgs,
		DisableFlagParsing: true,
	}
	cmd.AddCommand(createCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var (
		model  string
		prompt string
	)
	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), model, prompt)
		},
	}
	cmd.Flags().StringVar(&model, "model", "", "Model to be used")
	cmd.Flags().StringVar(&prompt, "prompt", "", "prompt")
	_ = cmd.MarkFlagRequired("model")
	_ = cmd.MarkFlagRequired("prompt")
	return cmd
}

func create(
	ctx context.Context,
	model string,
	prompt string,
) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := iv1.CreateCompletionRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	}
	body, err := ihttp.NewClient(env).SendRequest(http.MethodPost, path, &req)
	if err != nil {
		return err
	}
	scanner := sse.NewScanner(body)

	for scanner.Scan() {
		resp := scanner.Text()
		if !strings.HasPrefix(resp, "data: ") {
			// TODO(kenji): Handle other case.
			continue
		}

		respD := resp[5:]
		if respD == " [DONE]" {
			break
		}

		var d iv1.Completion
		if err := json.Unmarshal([]byte(respD), &d); err != nil {
			return fmt.Errorf("unmarshal response: %s", err)
		}
		cs := d.Choices
		if len(cs) == 0 {
			return fmt.Errorf("no choices")
		}
		fmt.Printf(cs[0].Text)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Printf("\n")

	return nil
}
