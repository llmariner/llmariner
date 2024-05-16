package completions

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	ihttp "github.com/llm-operator/cli/internal/http"
	"github.com/llm-operator/cli/internal/runtime"
	iv1 "github.com/llm-operator/inference-manager/api/v1"
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
		DisableFlagParsing: false,
	}
	cmd.AddCommand(createCmd())
	return cmd
}

func createCmd() *cobra.Command {
	var (
		model   string
		role    string
		content string
	)
	cmd := &cobra.Command{
		Use:  "create",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd.Context(), model, role, content)
		},
	}
	cmd.Flags().StringVar(&model, "model", "", "Model to be used")
	cmd.Flags().StringVar(&role, "role", "", "Chat completion role")
	cmd.Flags().StringVar(&content, "completion", "", "Chat completion content")
	_ = cmd.MarkFlagRequired("model")
	_ = cmd.MarkFlagRequired("role")
	_ = cmd.MarkFlagRequired("completion")
	return cmd
}

type delta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type choice struct {
	Delta delta `json:"delta"`
}

type data struct {
	Choices []choice `json:"choices"`
}

func create(
	ctx context.Context,
	model string,
	role string,
	content string,
) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}

	req := iv1.CreateChatCompletionRequest{
		Model: model,
		Messages: []*iv1.CreateChatCompletionRequest_Message{
			{
				Role:    role,
				Content: content,
			},
		},
		Stream: true,
	}
	body, err := ihttp.NewClient(env).SendRequest(http.MethodPost, path, &req)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 4096), 4096)
	scanner.Split(split)

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

		var d data
		if err := json.Unmarshal([]byte(respD), &d); err != nil {
			return fmt.Errorf("unmarshal response: %s", err)
		}
		cs := d.Choices
		if len(cs) == 0 {
			return fmt.Errorf("no choices")
		}
		fmt.Printf(cs[0].Delta.Content)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Printf("\n")

	return nil
}

// split tokenizes the input. The arguments are an initial substring of the remaining unprocessed data and a flag,
// atEOF, that reports whether the Reader has no more data to give.
// The return values are the number of bytes to advance the input and the next token to return to the user,
// if any, plus an error, if any.
func split(data []byte, atEOF bool) (int, []byte, error) {
	// Find a double newline.
	delims := [][]byte{
		[]byte("\r\r"),
		[]byte("\n\n"),
		[]byte("\r\n\r\n"),
	}
	pos := -1
	var dlen int
	for _, d := range delims {
		n := bytes.Index(data, d)
		if pos < 0 || (n >= 0 && n < pos) {
			pos = n
			dlen = len(d)
		}
	}

	if pos >= 0 {
		return pos + dlen, data[0:pos], nil
	}

	return 0, nil, nil
}
