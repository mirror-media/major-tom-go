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
	SlackAppToken string `yaml:"slackAppToken"`
	SlackBotToken string `yaml:"slackBotToken"`
}

type KubernetesConfigsRepo struct {
	Git     GitConfig  `yaml:"git"`
	Configs []Codebase `yaml:"configs"`
}

type Codebase struct {
	Projects []string `yaml:"projects"`
	Repo     string   `yaml:"repo"`
	Services []string `yaml:"services"`
	Stages   []string `yaml:"stages"`
	Type     int8     `yaml:"type"`
}

type Service struct {
	Name          string
	Repo          string
	SimpleService string
}

func contains(s []string, target string) bool {
	for _, stage := range s {
		if stage == target {
			return true
		}
	}
	return false
}

func (c Codebase) GetServices() (services []Service, err error) {
	switch c.Type {
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

func (c Codebase) getType1RepoPath(filename, stage string) (path string, err error) {
	if !contains(c.Stages, stage) {
		err = errors.Wrap(err, fmt.Sprintf("stage(%s) is not supported for %s", stage, c.Repo))
	}
	return fmt.Sprintf("%s/overlays/%s/%s", c.Repo, stage, filename), err
}

func (c Codebase) getType2ProjectPath(filename, stage, project string) (path string, err error) {
	if !contains(c.Stages, stage) {
		err = errors.Wrap(err, fmt.Sprintf("stage(%s) is not supported for %s", stage, c.Repo))
	}
	if !contains(c.Projects, project) {
		err = errors.Wrap(err, fmt.Sprintf("project(%s) is not supported for %s", project, c.Repo))
	}
	return fmt.Sprintf("%s/overlays/%s/overlays/%s/base/%s", c.Repo, stage, project, filename), err
}
func (c Codebase) getType2StagePath(filename, stage string) (path string, err error) {
	if !contains(c.Stages, stage) {
		err = errors.Wrap(err, fmt.Sprintf("stage(%s) is not supported for %s", stage, c.Repo))
	}
	return fmt.Sprintf("%s/overlays/%s/base/%s", c.Repo, stage, filename), err
}

func (c Codebase) getType2ServicePath(filename, stage, project, service string) (path string, err error) {
	if !contains(c.Stages, stage) {
		err = errors.Wrap(err, fmt.Sprintf("stage(%s) is not supported for %s", stage, c.Repo))
	}
	if !contains(c.Projects, project) {
		err = errors.Wrap(err, fmt.Sprintf("project(%s) is not supported for %s", project, c.Repo))
	}
	if !contains(c.Services, service) {
		err = errors.Wrap(err, fmt.Sprintf("service(%s) is not supported for %s", service, c.Repo))
	}
	return fmt.Sprintf("%s/overlays/%s/overlays/%s/overlays/%s/%s", c.Repo, stage, project, service, filename), err
}

func (c Codebase) GetImageKustomizationPath(stage, project string) (path string, err error) {
	switch c.Type {
	case 1:
		path, err = c.getType1RepoPath("kustomization.yaml", stage)
	case 2:
		if stage == "prod" {
			path, err = c.getType2ProjectPath("kustomization.yaml", stage, project)
		} else {
			path, err = c.getType2StagePath("kustomization.yaml", stage)
		}
	default:
		err = errors.New("StructureType is not supported. Check kubernetes-configs.yaml of major-tom-go")
	}
	return path, err
}

func (c Codebase) GetHpaPath(stage, project, service string) (path string, err error) {
	switch c.Type {
	case 1:
		path, err = c.getType1RepoPath("hpa.yaml", stage)
	case 2:
		path, err = c.getType2ServicePath("hpa.yaml", stage, project, service)
	default:
		err = errors.New("StructureType is not supported. Check kubernetes-configs.yaml of major-tom-go")
	}
	return path, err
}
