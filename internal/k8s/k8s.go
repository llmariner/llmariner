package k8s

import (
	"fmt"
	"time"

	"github.com/llm-operator/cli/internal/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewClient creates a new Kubernetes client.
func NewClient(env *runtime.Env) (kubernetes.Interface, error) {
	return kubernetes.NewForConfig(NewConfig(env))
}

// NewConfig creates a new Kubernetes configuration.
func NewConfig(env *runtime.Env) *rest.Config {
	return &rest.Config{
		Host:        fmt.Sprintf("%s/sessions", env.Config.EndpointURL),
		BearerToken: env.Token.AccessToken,
		Timeout:     30 * time.Second,
	}
}
