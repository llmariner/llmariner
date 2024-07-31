package k8s

import (
	"context"
	"io"
	"os"

	"github.com/llm-operator/cli/internal/runtime"
	corev1 "k8s.io/api/core/v1"
)

// StreamPodLogs streams logs from the given pod to stdout.
func StreamPodLogs(ctx context.Context, clusterID string, pod *corev1.Pod, follow bool) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	kc, err := NewClient(env, clusterID)
	if err != nil {
		return err
	}

	req := kc.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
		Follow: follow,
	})
	stream, err := req.Stream(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = stream.Close()
	}()
	_, err = io.Copy(os.Stdout, stream)
	return err
}
