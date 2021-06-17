package k8sop

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	gootkitconfig "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

type ReleaseVersion struct {
	AutoScaling struct {
		Enabled     bool
		MaxReplicas int
		MinReplicas int
	}
	ImageTag     string
	ReplicaCount int
	Version      int
}

func ListHelmReleaseVersion(kubeConfigPath string, name, namespace string, maxHistory int) ([]ReleaseVersion, error) {
	if maxHistory == 0 {
		maxHistory = 1
	}

	history, err := getHelmReleaseHistory(kubeConfigPath, name, namespace, 11)
	if err != nil {
		return nil, err
	}

	sort.Slice(history, func(i, j int) bool { return history[i].Version > history[j].Version })

	releaseVersions := make([]ReleaseVersion, 0, len(history))

	for _, version := range history {
		autoScaling := version.Config["autoscaling"].(map[string]interface{})

		releaseVersions = append(releaseVersions, ReleaseVersion{
			ImageTag:     version.Config["image"].(map[string]interface{})["tag"].(string),
			ReplicaCount: int(version.Config["replicaCount"].(float64)),
			Version:      version.Version,
			AutoScaling: struct {
				Enabled     bool
				MaxReplicas int
				MinReplicas int
			}{
				Enabled:     autoScaling["enabled"].(bool),
				MaxReplicas: int(autoScaling["maxReplicas"].(float64)),
				MinReplicas: int(autoScaling["minReplicas"].(float64)),
			},
		})
	}
	return releaseVersions, err
}

func getHelmRelease(kubeConfigPath, name, namespace string, version int) (*release.Release, error) {

	actionConfig := new(action.Configuration)

	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(kube.GetConfig(kubeConfigPath, "", namespace), namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		logrus.Warn(err)
		return nil, err
	}

	get := action.NewGet(actionConfig)

	get.Version = version
	release, err := get.Run(name)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	return release, nil
}

func getHelmReleaseHistory(kubeConfigPath, name, namespace string, max int) ([]*release.Release, error) {

	actionConfig := new(action.Configuration)

	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(kube.GetConfig(kubeConfigPath, "", namespace), namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		logrus.Warn(err)
		return nil, err
	}

	history := action.NewHistory(actionConfig)
	history.Max = max

	releases, err := history.Run(name)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	return releases, nil
}

func parseHelmReleaseVersion(name, manifest string, version int) (ReleaseVersion, error) {

	values := strings.Split(manifest, "COMPUTED VALUES:\n")[1]
	computed_values := strings.Split(values, "---\n")[0]

	fmt.Println(computed_values)
	parser := gootkitconfig.New("release" + ":" + name + ":" + strconv.Itoa(version))
	defer parser.ClearAll()

	parser.AddDriver(yaml.Driver)
	err := parser.LoadStrings(gootkitconfig.Yaml, computed_values)
	if err != nil {
		logrus.Warn(err)
		return ReleaseVersion{}, err
	}

	return ReleaseVersion{
		ImageTag:     parser.String("image.tag", ""),
		ReplicaCount: parser.Int("replicaCount", -1),
		Version:      version,
		AutoScaling: struct {
			Enabled     bool
			MaxReplicas int
			MinReplicas int
		}{
			Enabled:     parser.Bool("autoscaling.enabled", false),
			MaxReplicas: parser.Int("autoscaling.maxReplicas", -1),
			MinReplicas: parser.Int("autoscaling.minReplicas", -1),
		},
	}, nil
}
