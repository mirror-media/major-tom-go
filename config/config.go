package config

type KubeConfigPath string
type Stage string
type Project string
type K8S map[Project]map[Stage]KubeConfigPath

type Config struct {
	SlackBotToken  string `yaml:"slackBotToken"`
	SlackToken     string `yaml:"slackToken"`
	ClusterConfigs K8S    `yaml:"clusterConfigs"`
}
