package k8sop

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
// kubeConfigPath = "/dummypath"
)

func getPodInfo(ctx context.Context, kubeConfigPath string, namespace string, name string) (map[string]int, error) {

	// Initialize kubernetes-client
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	// Create new client with the given config
	// https://pkg.go.dev/k8s.io/client-go/kubernetes?tab=doc#NewForConfig
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Use the app's label selector name. Remember this should match with
	// the deployment selector's matchLabels.
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=" + name,
	})
	if err != nil {
		return nil, err
	}
	status := make(map[string]int)
	for _, pod := range list.Items {
		var ready string
		for _, cond := range pod.Status.Conditions {
			if cond.Type == "Ready" {
				ready = string(cond.Status)
			}
		}
		imageParts := strings.Split(pod.Spec.Containers[0].Image, ":")
		key := fmt.Sprintf("%s, Phase: %s, Ready: %s", imageParts[len(imageParts)-1], pod.Status.Phase, ready)
		if _, found := status[key]; found {
			status[key] = 0
		}
		status[key]++
	}
	return status, nil
}
