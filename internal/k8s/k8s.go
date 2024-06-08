package k8s

import (
	"fmt"

	"github.com/llm-operator/cli/internal/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewClient creates a new Kubernetes client.
func NewClient(env *runtime.Env, clusterID string) (kubernetes.Interface, error) {
	return kubernetes.NewForConfig(newConfig(env, clusterID))
}

// newConfig creates a new Kubernetes configuration.
func newConfig(env *runtime.Env, clusterID string) *rest.Config {
	return &rest.Config{
		Host:        fmt.Sprintf("%s/sessions/%s", env.Config.EndpointURL, clusterID),
		BearerToken: env.Token.AccessToken,
	}
}
