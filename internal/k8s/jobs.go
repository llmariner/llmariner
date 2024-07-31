package k8s

import (
	"context"
	"fmt"

	"github.com/llm-operator/cli/internal/runtime"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListPodsForJob lists all pods for the given job.
func ListPodsForJob(ctx context.Context, clusterID, ns, name string) ([]corev1.Pod, error) {
	env, err := runtime.NewEnv(ctx)
	if err != nil {
		return nil, err
	}
	kc, err := NewClient(env, clusterID)
	if err != nil {
		return nil, err
	}

	resp, err := kc.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
		// This is an implicit assumption that the job name is equal to "<job_id>".
		LabelSelector: fmt.Sprintf("job-name=%s", name),
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("no pod found for the job %q", name)
	}
	return resp.Items, nil
}

// FindLatestOrLastFailedPod finds the latest pod and the last failed pod from the given list of pods.
func FindLatestOrLastFailedPod(pods []corev1.Pod) (latest *corev1.Pod, lastFailed *corev1.Pod) {
	for _, pod := range pods {
		if latest == nil || pod.CreationTimestamp.After(latest.CreationTimestamp.Time) {
			latest = &pod
		}

		if pod.Status.Phase != corev1.PodFailed {
			continue
		}
		if lastFailed == nil || pod.CreationTimestamp.After(lastFailed.CreationTimestamp.Time) {
			lastFailed = &pod
		}
	}
	return latest, lastFailed
}
