package k8s

import (
	"time"

	"github.com/llm-operator/cli/internal/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewClient creates a new Kubernetes client.
func NewClient(env *runtime.Env) (kubernetes.Interface, error) {
	// Remove "/v1" from the endpoint URL. This is a hack to construct the proper URL for k8s.
	host := env.Config.EndpointURL[:len(env.Config.EndpointURL)-2]
	return kubernetes.NewForConfig(&rest.Config{
		Host:        host,
		BearerToken: env.Token.AccessToken,
		Timeout:     30 * time.Second,
	})
}
