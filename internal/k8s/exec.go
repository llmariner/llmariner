package k8s

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/docker/cli/cli/streams"
	"github.com/llmariner/llmariner/internal/runtime"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/httpstream"
	rc "k8s.io/apimachinery/pkg/util/remotecommand"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecPod exects into a pod.
func ExecPod(ctx context.Context, clusterID string, pod *corev1.Pod) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil
	}
	kc, err := NewClient(env, clusterID)
	if err != nil {
		return err
	}

	req := kc.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec").
		VersionedParams(
			&corev1.PodExecOptions{
				Command: []string{"/bin/bash"},
				Stdin:   true,
				Stdout:  true,
				Stderr:  true,
				TTY:     true,
			},
			scheme.ParameterCodec,
		)

	config := newConfig(env, clusterID)

	// Use SPDY or websocket.
	// Please note that Kong Ingress Controller does not support SPDY (https://github.com/Kong/kong/discussions/7334).
	spdyExec, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
	if err != nil {
		return err
	}
	// V5 is the latest protocol that is used from NewWebSocketExecutor(), but a k8s cluster may not support it.
	wsExec, err := remotecommand.NewWebSocketExecutorForProtocols(config, http.MethodGet, req.URL().String(),
		rc.StreamProtocolV5Name,
		rc.StreamProtocolV4Name,
		rc.StreamProtocolV3Name,
	)
	if err != nil {
		return err
	}

	exec, err := remotecommand.NewFallbackExecutor(wsExec, spdyExec, httpstream.IsUpgradeFailure)
	if err != nil {
		return err
	}

	in := streams.NewIn(os.Stdin)
	if err := in.SetRawTerminal(); err != nil {
		return err
	}
	defer in.RestoreTerminal()

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  in,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return fmt.Errorf("stream: %s", err)
	}

	return nil
}
