package completions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	iv1 "github.com/llmariner/inference-manager/api/v1"
	"github.com/llmariner/inference-manager/common/pkg/sse"
	ihttp "github.com/llmariner/llmariner/cli/internal/http"
	"github.com/llmariner/llmariner/cli/internal/runtime"
	"github.com/spf13/cobra"
)

const (
	path = "/chat/completions"
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
	msg := &iv1.CreateChatCompletionRequest_Message{}
	req := &iv1.CreateChatCompletionRequest{
		Messages: []*iv1.CreateChatCompletionRequest_Message{msg},
		Stream:   true,
	}

	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), req)
		},
	}
	cmd.Flags().StringVar(&req.Model, "model", "", "Model to be used")
	cmd.Flags().StringVar(&msg.Role, "role", "", "Chat completion role")
	cmd.Flags().StringVar(&msg.Content, "completion", "", "Chat completion content")
	cmd.Flags().Float64Var(&req.PresencePenalty, "presence-penalty", 0.0, "Presence penalty")
	cmd.Flags().Float64Var(&req.FrequencyPenalty, "frequency-penalty", 0.0, "Frequency penalty")
	cmd.Flags().StringArrayVar(&req.Stop, "stop", nil, "Stop words")
	cmd.Flags().Int32Var(&req.MaxTokens, "max-tokens", 0, "Max tokens")
	cmd.Flags().Float64Var(&req.Temperature, "temperature", 1.0, "Temperature")
	cmd.Flags().Float64Var(&req.TopP, "top-p", 1.0, "Top p")
	_ = cmd.MarkFlagRequired("model")
	_ = cmd.MarkFlagRequired("role")
	_ = cmd.MarkFlagRequired("completion")
	return cmd
}

func create(ctx context.Context, req *iv1.CreateChatCompletionRequest) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
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

		var d iv1.ChatCompletionChunk
		if err := json.Unmarshal([]byte(respD), &d); err != nil {
			return fmt.Errorf("unmarshal response: %s", err)
		}
		cs := d.Choices
		if len(cs) == 0 {
			return fmt.Errorf("no choices")
		}
		fmt.Print(cs[0].Delta.Content)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Print("\n")

	return nil
}
