package k8s

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/llmariner/llmariner/internal/runtime"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// PortForward establishes a port-forwarding session to a pod.
//
// This first attempts tunneling (websocket) dialer, then fallback to SPDY dialer.
// Please note that the websocket works only with k8s 1.30 with the "PortForwardWebsockets"
// feature gate enabled.
func PortForward(ctx context.Context, clusterID string, pod *corev1.Pod, localPort, remotePort int) error {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil
	}

	config := newConfig(env, clusterID)
	url, err := url.Parse(config.Host)
	if err != nil {
		return err
	}

	url.Path = path.Join(
		"v1",
		"sessions",
		"api",
		"v1",
		"namespaces",
		pod.Namespace,
		"pods",
		pod.Name,
		"portforward",
	)

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}
	spdyDialer := spdy.NewDialer(
		upgrader,
		&http.Client{Transport: transport},
		http.MethodPost,
		url,
	)

	tunnelingDialer, err := portforward.NewSPDYOverWebsocketDialer(url, config)
	if err != nil {
		return err
	}
	dialer := portforward.NewFallbackDialer(tunnelingDialer, spdyDialer, httpstream.IsUpgradeFailure)

	stopCh := make(chan struct{}, 1)
	readyCh := make(chan struct{})

	fw, err := portforward.New(
		dialer,
		[]string{fmt.Sprintf("%d:%d", localPort, remotePort)},
		stopCh,
		readyCh,
		os.Stdout,
		os.Stderr,
	)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}
