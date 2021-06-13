package test

import "github.com/mirror-media/major-tom-go/config"

// FIXME we need proper test config
var GitConfigsTest = map[config.Repository]config.GitConfig{
	"tv": config.GitConfig{
		Branch:        "test/majortom",
		SSHKeyPath:    "/Users/chiu/dev/mtv/major-tom-go/configs/ssh/identity",
		SSHKeyUser:    "mnews@mnews.tw",
		SSHKnownhosts: "/Users/chiu/dev/mtv/major-tom-go/configs/ssh/known_hosts",
		URL:           "ssh://source.developers.google.com:2022/p/mirror-tv-275709/r/helm",
	},
}

var ConfigTest = config.Config{
	ClusterConfigs: config.K8S{
		// FIXME we need proper test path
		// "mm": {
		// 	"prod":    "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		// 	"staging": "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		// 	"dev":     "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		// },
		"tv": {
			"prod":    "/Users/chiu/dev/mtv/major-tom-go/configs/kubeconfig/kubeconfig-prod-tv",
			"staging": "/Users/chiu/dev/mtv/major-tom-go/configs/kubeconfig/kubeconfig-staging-tv",
			"dev":     "/Users/chiu/dev/mtv/major-tom-go/configs/kubeconfig/kubeconfig-dev-tv",
		},
		// "readr": {
		// 	"prod": "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		// 	"dev":  "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		// },
	},
	GitConfigs: GitConfigsTest,
}
