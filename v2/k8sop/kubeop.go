package k8sop

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
// kubeConfigPath = "/dummypath"
)

type DeploymentInfo struct {
	Available int32
	ImageTag  string
	Ready     int32
	Updated   int32
}

func getKubeCliSet(kubeConfigPath string, namespace string) (clientset *kubernetes.Clientset, err error) {
	// Initialize kubernetes-client
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	// Create new client with the given config
	// https://pkg.go.dev/k8s.io/client-go/kubernetes?tab=doc#NewForConfig
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, err
}

func getKubeDynamicCliSet(kubeConfigPath string, namespace string) (clientset dynamic.Interface, err error) {
	// Initialize kubernetes-client
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	// Create new client with the given config
	// https://pkg.go.dev/k8s.io/client-go/kubernetes?tab=doc#NewForConfig
	clientset, err = dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, err
}

type ReleaseInfo struct {
	Name      string `json:"releaseName"`
	Namespace string
	Status    string `json:"releaseStatus"`
}

type Metadata struct {
	Namespace string `json:"namespace"`
}

type ReleaseStatus struct {
	Status   ReleaseInfo `json:"status"`
	Metadata Metadata    `json:"metadata"`
}

func ListReleases(ctx context.Context, kubeConfigPath string) (releaseInfo []ReleaseInfo, err error) {
	return listReleases(ctx, kubeConfigPath, "default")
}

func listReleases(ctx context.Context, kubeConfigPath string, namespace string) (releaseInfo []ReleaseInfo, err error) {
	clientset, err := getKubeDynamicCliSet(kubeConfigPath, namespace)
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    "helm.fluxcd.io",
		Version:  "v1",
		Resource: "helmreleases",
	}

	list, err := clientset.Resource(gvr).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := list.Items

	releaseInfo = make([]ReleaseInfo, len(items))
	for i, item := range items {
		var status ReleaseStatus

		b, err := item.MarshalJSON()
		if err != nil {
			return nil, err
		}
		json.Unmarshal(b, &status)
		releaseInfo[i] = status.Status
		releaseInfo[i].Namespace = status.Metadata.Namespace
	}

	sort.Slice(releaseInfo, func(i, j int) bool { return releaseInfo[i].Name < releaseInfo[j].Name })

	return releaseInfo, nil
}

// GetDeploymentInfo return the status of current deployments for the specific service
func GetDeploymentInfo(ctx context.Context, kubeConfigPath string, name string) (DeploymentInfo, error) {
	list, err := ListReleases(ctx, kubeConfigPath)
	if err != nil {
		return DeploymentInfo{}, err
	}
	namespaces := make(map[string]string, len(list))
	for _, release := range list {
		namespaces[release.Name] = release.Namespace
	}
	namespace, isExisting := namespaces[name]
	if !isExisting {
		return DeploymentInfo{}, errors.Errorf("service(%s) doesn't exist or isn't operable", name)
	}
	return getDeploymentInfo(ctx, kubeConfigPath, namespace, name)
}

func getDeploymentInfo(ctx context.Context, kubeConfigPath string, namespace string, name string) (DeploymentInfo, error) {

	clientset, err := getKubeCliSet(kubeConfigPath, namespace)
	if err != nil {
		return DeploymentInfo{}, err
	}

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, v1.GetOptions{
		TypeMeta: v1.TypeMeta{
			Kind: "Deployment",
		},
	})
	if err != nil {
		return DeploymentInfo{}, err
	}

	containers := deployment.Spec.Template.Spec.Containers
	imageParts := strings.Split(containers[0].Image, ":")

	return DeploymentInfo{
		Available: deployment.Status.AvailableReplicas,
		ImageTag:  imageParts[len(imageParts)-1],
		Ready:     deployment.Status.ReadyReplicas,
		Updated:   deployment.Status.UpdatedReplicas,
	}, nil

}

func getPodInfo(ctx context.Context, kubeConfigPath string, namespace string, name string) (map[string]int, error) {

	clientset, err := getKubeCliSet(kubeConfigPath, namespace)
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
