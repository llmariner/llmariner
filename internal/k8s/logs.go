package k8s

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/llmariner/cli/internal/runtime"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

var colors = []color.Attribute{
	color.FgBlue,
	color.FgGreen,
	color.FgRed,
	color.FgCyan,
	color.FgYellow,
	color.FgMagenta,
	color.FgHiBlue,
	color.FgHiGreen,
	color.FgHiRed,
	color.FgHiCyan,
	color.FgHiYellow,
	color.FgHiMagenta,
}

// StreamPodsLogs streams logs from the given pod to stdout.
func StreamPodsLogs(ctx context.Context, clusterID string, follow bool, pods ...corev1.Pod) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return err
	}
	kc, err := NewClient(env, clusterID)
	if err != nil {
		return err
	}

	enablePrefix := len(pods) > 1
	g, ctx := errgroup.WithContext(ctx)
	for i, pod := range pods {
		var prefix string
		if enablePrefix {
			c := color.New(colors[i%len(colors)]).SprintFunc()
			prefix = c(fmt.Sprintf("%s: ", pod.Name))
		}
		g.Go(func() error { return streamPodLogs(ctx, kc, follow, &pod, prefix) })
	}
	return g.Wait()
}

func streamPodLogs(ctx context.Context, client kubernetes.Interface, follow bool, pod *corev1.Pod, prefix string) error {
	req := client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Follow: follow})
	stream, err := req.Stream(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = stream.Close() }()

	r := bufio.NewReader(stream)
	for {
		line, err := r.ReadBytes('\n')
		if len(line) != 0 {
			_, _ = fmt.Fprintf(os.Stdout, "%s%s", prefix, line)
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
