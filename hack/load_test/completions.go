package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	iv1 "github.com/llmariner/inference-manager/api/v1"
	"github.com/llmariner/inference-manager/common/pkg/sse"
)

const (
	path = "/chat/completions"
)

func sendChatCompletion(
	ctx context.Context,
	endpointURL string,
	accessToken string,
	req *iv1.CreateChatCompletionRequest,
	printOutput bool,
	reqID string,
) error {
	client := newClient(endpointURL, accessToken)
	body, err := client.SendRequest(http.MethodPost, path, &req, reqID)
	if err != nil {
		return err
	}

	if !req.Stream {
		b, err := io.ReadAll(body)
		if err != nil {
			return fmt.Errorf("read response: %s", err)
		}

		// Parse the response as ChatCompletion.
		var resp iv1.ChatCompletion
		if err := json.Unmarshal(b, &resp); err != nil {
			return fmt.Errorf("unmarshal response: %s", err)
		}
		if printOutput {
			for _, c := range resp.Choices {
				fmt.Print(c.Message.Content)
			}
			fmt.Print("\n")
		}
		return nil
	}

	// Process the streaming response.

	scanner := sse.NewScanner(body)

	var respBuf string

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
			fmt.Println("Error unmarshalling response:", err)
			fmt.Printf("ResponseBuf: %s\n", respBuf)
			fmt.Printf("Response: %s\n", respD)
			return fmt.Errorf("unmarshal response: %s", err)
		}
		cs := d.Choices
		if len(cs) > 0 && printOutput {
			// TODO(kenji): Handle multiple choices.
			fmt.Print(cs[0].Delta.Content)
			respBuf += cs[0].Delta.Content
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	if printOutput {
		fmt.Print("\n")
	}

	return nil
}
