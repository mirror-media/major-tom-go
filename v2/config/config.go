package config

import (
	"fmt"
	"sort"

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
	Project       string
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
		for _, project := range c.Projects {
			for _, service := range c.Services {
				services = append(services, Service{
					Name:          fmt.Sprintf("%s-%s-%s", c.Repo, project, service),
					Repo:          c.Repo,
					SimpleService: service,
				})
			}
		}
		sort.Slice(services, func(i, j int) bool {
			return services[i].Name < services[j].Name
		})
	default:
		err = errors.New("Type is not supported. Check kubernetes-configs.yaml of major-tom-go")
	}
	return services, err
}

func (c Codebase) getType1StagePath(filename, stage string) (path string, err error) {
	path = fmt.Sprintf("%s/overlays/%s/%s", c.Repo, stage, filename)
	if c.Type != 1 {
		return path, errors.New(fmt.Sprintf("codebase has type(%d) so the path is wrong", c.Type))
	}
	if !contains(c.Stages, stage) {
		return path, errors.New(fmt.Sprintf("stage(%s) is not supported for %s", stage, c.Repo))
	}
	return path, err
}

func (c Codebase) getType2ProjectPath(filename, stage, project string) (path string, err error) {
	path = fmt.Sprintf("%s/overlays/%s/overlays/%s/base/%s", c.Repo, stage, project, filename)

	if c.Type != 2 {
		return path, errors.New(fmt.Sprintf("codebase has type(%d) so the path is wrong", c.Type))
	}

	if !contains(c.Stages, stage) {
		return path, errors.New(fmt.Sprintf("stage(%s) is not supported for %s", stage, c.Repo))
	}
	if !contains(c.Projects, project) {
		return path, errors.New(fmt.Sprintf("project(%s) is not supported for %s", project, c.Repo))
	}
	return path, err
}

func (c Codebase) getType2StagePath(filename, stage string) (path string, err error) {
	path = fmt.Sprintf("%s/overlays/%s/base/%s", c.Repo, stage, filename)

	if c.Type != 2 {
		return path, errors.New(fmt.Sprintf("codebase has type(%d) so the path is wrong", c.Type))
	}
	if !contains(c.Stages, stage) {
		return path, errors.New(fmt.Sprintf("stage(%s) is not supported for %s", stage, c.Repo))
	}
	return path, err
}

func (c Codebase) getType2ServicePath(filename, stage, project, service string) (path string, err error) {
	path = fmt.Sprintf("%s/overlays/%s/overlays/%s/overlays/%s/%s", c.Repo, stage, project, service, filename)

	if c.Type != 2 {
		return path, errors.New(fmt.Sprintf("codebase has type(%d) so the path is wrong", c.Type))
	}
	if !contains(c.Stages, stage) {
		return path, errors.New(fmt.Sprintf("stage(%s) is not supported for %s", stage, c.Repo))
	}
	if !contains(c.Projects, project) {
		return path, errors.New(fmt.Sprintf("project(%s) is not supported for %s", project, c.Repo))
	}
	if !contains(c.Services, service) {
		return path, errors.New(fmt.Sprintf("service(%s) is not supported for %s", service, c.Repo))
	}
	return path, err
}

func (c Codebase) GetImageKustomizationPath(stage, project string) (path string, err error) {
	switch c.Type {
	case 1:
		path, err = c.getType1StagePath("kustomization.yaml", stage)
	case 2:
		if stage == "prod" {
			path, err = c.getType2ProjectPath("kustomization.yaml", stage, project)
		} else {
			path, err = c.getType2StagePath("kustomization.yaml", stage)
		}
	default:
		err = errors.New("Type is not supported. Check kubernetes-configs.yaml of major-tom-go")
	}
	return path, err
}

func (c Codebase) GetHpaPath(stage, project, service string) (path string, err error) {
	switch c.Type {
	case 1:
		path, err = c.getType1StagePath("hpa.yaml", stage)
	case 2:
		path, err = c.getType2ServicePath("hpa.yaml", stage, project, service)
	default:
		err = errors.New("Type is not supported. Check kubernetes-configs.yaml of major-tom-go")
	}
	return path, err
}
