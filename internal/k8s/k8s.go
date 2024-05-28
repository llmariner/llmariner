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
	host := fmt.Sprintf("%s/sessions", env.Config.EndpointURL)
	return kubernetes.NewForConfig(&rest.Config{
		Host:        host,
		BearerToken: env.Token.AccessToken,
		Timeout:     30 * time.Second,
	})
}
