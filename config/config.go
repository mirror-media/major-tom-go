package config

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
	SlackBotToken  string                   `yaml:"slackBotToken"`
	SlackToken     string                   `yaml:"slackToken"`
	ClusterConfigs K8S                      `yaml:"clusterConfigs"`
	GitConfigs     map[Repository]GitConfig `yaml:"gitConfigs"`
}
