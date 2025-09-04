package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	iv1 "github.com/llmariner/inference-manager/api/v1"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	parallelism  = 1
	duration     = 5 * time.Minute
	logFrequency = 1

	model             = "ft:TinyLlama-TinyLlama-1.1B-Chat-v1.0:test"
	endpointURL       = "http://localhost:8080/v1"
	accessTokenEnvVar = "LLMARINER_API_KEY"

	llmarinerNamespace     = "llmariner"
	kubeconfigRelativePath = ".kube/config"

	printOutput = false
)

func injectFaults(ctx context.Context) error {
	kubeconfigPath := os.Getenv("HOME") + "/" + kubeconfigRelativePath
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	for {
		apps := []string{
			"inference-manager-engine",
			"inference-manager-server",
		}

		var pods []corev1.Pod
		for _, app := range apps {
			resp, err := clientset.CoreV1().Pods(llmarinerNamespace).List(ctx, metav1.ListOptions{
				LabelSelector: fmt.Sprintf("app.kubernetes.io/name=%s", app),
			})
			if err != nil {
				return err
			}

			for _, pod := range resp.Items {
				pods = append(pods, pod)
			}
		}

		if len(pods) == 0 {
			fmt.Printf("No pods found in namespace %s. Retrying...\n", llmarinerNamespace)
			time.Sleep(10 * time.Second)
			continue
		}

		// Randomly select a pod to delete.
		pod := pods[rand.Intn(len(pods))]
		fmt.Printf("Deleting pod %q\n", pod.Name)
		if err := clientset.CoreV1().Pods(llmarinerNamespace).Delete(ctx, pod.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}

		// Wait for a while before injecting the next fault.
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(90 * time.Second):
		}
	}
}

func runLoadTest() error {
	start := time.Now()
	accessToken := os.Getenv(accessTokenEnvVar)
	if accessToken == "" {
		return fmt.Errorf("environment variable %s is not set", accessTokenEnvVar)
	}

	messages := []string{
		"Hello, World!",
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	count := &atomic.Int64{}
	errCount := &atomic.Int64{}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return injectFaults(ctx)
	})

	for i := 0; i < parallelism; i++ {
		i := i
		eg.Go(func() error {
			j := 0
			for {
				j++
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

				reqID := fmt.Sprintf("load-test-%d-%d", i, j)

				if err := sendChatCompletion(ctx, endpointURL, accessToken, req, printOutput, reqID); err != nil {
					fmt.Printf("Error sending request, reqID=%s, (%s): %s\n", reqID, time.Now(), err)
					errCount.Add(1)
					continue
				}

				n := count.Add(1)

				if n%logFrequency == 0 {
					fmt.Printf("Processed %d requests (%d errors). %s elapsed\n", n, errCount.Load(), time.Since(start))
				}

				time.Sleep(1 * time.Second)

				select {
				case <-ctx.Done():
					return nil
				default:
				}

			}
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error during load test: %s", err)
	}
	fmt.Printf("Load test completed. Total requests: %d, Errors: %d\n", count.Load(), errCount.Load())
	return nil
}

func main() {
	if err := runLoadTest(); err != nil {
		fmt.Printf("Error running load test: %s\n", err)
	}
}
