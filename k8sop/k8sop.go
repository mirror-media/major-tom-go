// Package k8sop is responsible of the implementation involving helm and Kubernetes
package k8sop

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const (
	DeploymentTagKey = "tag"
	kubeConfigPath   = "/dummypath"
)

func getHelmReleaseManifest(name string, namespace string) (string, error) {

	actionConfig := new(action.Configuration)

	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(kube.GetConfig(kubeConfigPath, "", namespace), namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		return "", err
	}

	get := action.NewGet(actionConfig)

	get.Version = 0
	release, err := get.Run(name)
	if err != nil {
		return "", err
	}
	return release.Manifest, nil
}

func GetHelmReleaseImageTag(name string) (map[string]string, error) {

	manifest, err := getHelmReleaseManifest(name, "default")
	if err != nil {
		return nil, err
	}

	type Deployment struct {
		Spec struct {
			Template struct {
				Spec struct {
					Containers []struct {
						Image string `yaml:"image"`
					} `yaml:"containers"`
				} `yaml:"spec"`
			} `yaml:"template"`
		} `yaml:"spec"`
	}

	resources := strings.Split(manifest, "---\n")
	var d Deployment
	for _, resource := range resources {
		if strings.Contains(resource, "kind: Deployment") {
			fmt.Println(resource)
			err = yaml.Unmarshal([]byte(resource), &d)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	imageParts := strings.Split(d.Spec.Template.Spec.Containers[0].Image, ":")
	tag := imageParts[len(imageParts)-1]

	return map[string]string{
		DeploymentTagKey: tag,
	}, nil
}
