package main

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	iv1 "github.com/llmariner/inference-manager/api/v1"
	"golang.org/x/sync/errgroup"
)

const (
	parallelism  = 50
	duration     = 5 * time.Minute
	logFrequency = 1000

	model             = "meta-llama-Llama-3.2-1B-Instruct"
	endpointURL       = "https://api.cloudnatix.com/v1"
	accessTokenEnvVar = "LLMARINER_TOKEN"

	printOutput = false
)

func runLoadTest() error {
	accessToken := os.Getenv(accessTokenEnvVar)
	if accessToken == "" {
		return fmt.Errorf("environment variable %s is not set", accessTokenEnvVar)
	}

	messages := []string{
		"Hello, World!",
		"Where is the capital of France?",
		"What is the meaning of life?",
		"Tell me a joke.",
		"How do you make a sandwich?",
		"What's the weather like today?",
		"Can you recommend a good book?",
		"What's your favorite movie?",
		"How do I cook pasta?",
		"Tell me a fun fact.",
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	count := &atomic.Int64{}

	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < parallelism; i++ {
		i := i
		eg.Go(func() error {
			for {
				message := messages[i%len(messages)]

				req := &iv1.CreateChatCompletionRequest{
					Model: model,
					Messages: []*iv1.CreateChatCompletionRequest_Message{
						{
							Role: "user",
							Content: []*iv1.CreateChatCompletionRequest_Message_Content{
								{
									Type: "text",
									Text: message,
								},
							},
						},
					},
					Stream: true,
				}

				if err := sendChatCompletion(ctx, endpointURL, accessToken, req, printOutput); err != nil {
					return fmt.Errorf("send chat completion request: %w", err)
				}

				n := count.Add(1)

				if n%logFrequency == 0 {
					fmt.Printf("Processed %d requests\n", n)
				}

				select {
				case <-ctx.Done():
					return nil
				default:
				}

			}
		})
	}

	return eg.Wait()
}

func main() {
	if err := runLoadTest(); err != nil {
		fmt.Printf("Error running load test: %v\n", err)
	}
	fmt.Println("Load test completed successfully.")
}
