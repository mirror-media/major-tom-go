package config

import (
	"fmt"

	"github.com/pkg/errors"
)

type KubeConfigPath string
type Stage string
type Project string
type K8S map[Project]map[Stage]KubeConfigPath

type Repository string
type GitConfig struct {
	Branch        string `yaml:"branch"`
	SSHKeyPath    string `yaml:"sshKeyPath"`
	SSHKeyUser    string `yaml:"sshKeyUser"`
	SSHKnownhosts string `yaml:"sshKnownhosts"`
	URL           string `yaml:"url"`
}

type Config struct {
	ClusterConfigs K8S                      `yaml:"clusterConfigs"`
	GitConfigs     map[Repository]GitConfig `yaml:"gitConfigs"`
	SlackAppToken  string                   `yaml:"slackAppToken"`
	SlackBotToken  string                   `yaml:"slackBotToken"`
}

type KubernetesConfigs struct {
	Repo    string     `json:"repo"`
	Branch  string     `json:"branch"`
	Configs []Codebase `json:"configs"`
}

type Codebase struct {
	Projects      []string `json:"projects"`
	Repo          string   `json:"repo"`
	Services      []string `json:"services"`
	Stages        []string `json:"stages"`
	StructureType int64    `json:"structureType"`
}

type Service struct {
	Name          string
	Repo          string
	SimpleService string
}

func (c Codebase) GetServiceNames() (services []Service, err error) {
	switch c.StructureType {
	case 1:
		services = append(services, Service{
			Name: c.Repo,
			Repo: c.Repo,
		})
	case 2:
		for _, service := range c.Services {
			services = append(services, Service{
				Name:          fmt.Sprintf("%s-%s-%s", c.Repo, c.Projects, service),
				Repo:          c.Repo,
				SimpleService: service,
			})
		}
	default:
		err = errors.New("StructureType is not supported. Check kubernetes-configs.yaml of major-tom-go")
	}
	return services, err
}

func (c Codebase) getSimplePathByStage(filename, stage string) (path string, err error) {
	// TODO validate stage
	return fmt.Sprintf("%s/overlays/%s/base/%s", c.Repo, stage, filename), err
}
func (c Codebase) getComplexPathByStageAndProject(filename, stage, project string) (path string, err error) {
	// TODO validate stage and project
	return fmt.Sprintf("%s/overlays/%s/overlays/%s/base/%s", c.Repo, stage, project, filename), err
}

func (c Codebase) getComplexServicePathByStageAndProject(filename, stage, project, service string) (path string, err error) {
	return fmt.Sprintf("%s/overlays/%s/overlays/%s/overlays/%s/%s", c.Repo, stage, project, service, filename), err
}

func (c Codebase) GetImageKustomizationPathByStage(stage string) (path string, err error) {
	return c.getSimplePathByStage("kustomization.yaml", stage)
}

func (c Codebase) GetImageKustomizationForProdByProject(project string) (path string, err error) {
	switch c.StructureType {
	case 1:
		path, err = c.GetImageKustomizationPathByStage("prod")
	case 2:
		path, err = c.getComplexPathByStageAndProject("kustomization.yaml", c.Repo, project)
	default:
		err = errors.New("StructureType is not supported. Check kubernetes-configs.yaml of major-tom-go")
	}
	return path, err
}

func (c Codebase) GetSimpleHpaPathByStage(stage string) (path string, err error) {
	switch c.StructureType {
	case 1:
		path, err = c.getSimplePathByStage("hpa.yaml", "prod")
	case 2:
		err = errors.New("StructureType is 2. Use GetHpaPathForProdByProjectAndService() instead")
	default:
		err = errors.New("StructureType is not supported. Check kubernetes-configs.yaml of major-tom-go")
	}
	return path, err
}

func (c Codebase) GetHpaPathForProdByProjectAndService(project, service string) (path string, err error) {
	switch c.StructureType {
	case 1:
		err = errors.New("StructureType is 1. Use GetSimpleHpaPathByStage() instead")
	case 2:
		// TODO validate project
		path, err = c.getComplexServicePathByStageAndProject("hpa.yaml", "prod", project, service)
	default:
		err = errors.New("StructureType is not supported. Check kubernetes-configs.yaml of major-tom-go")
	}
	return path, err
}
