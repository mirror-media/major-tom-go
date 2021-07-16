package test

import "github.com/mirror-media/major-tom-go/v2/config"

var K8sRepo = config.KubernetesConfigsRepo{
	GitConfig: config.GitConfig{
		URL:           "git@github.com:mirror-media/kubernetes-configs.git",
		Branch:        "major-tom-test",
		SSHKeyPath:    "../configs/ssh/identity",
		SSHKeyUser:    "git",
		SSHKnownhosts: "../configs/ssh/known_hosts",
	},
	Configs: []config.Codebase{
		{
			Type:     2,
			Repo:     "openwarehouse",
			Stages:   []string{"dev", "staging", "prod"},
			Projects: []string{"tv"},
			Services: []string{"cms", "gql-external", "gql-internal"},
		},
		{
			Type:   1,
			Repo:   "mirror-tv-nuxt",
			Stages: []string{"dev", "staging", "prod"},
		},
	},
}
