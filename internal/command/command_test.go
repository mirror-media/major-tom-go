// Package command implement the operation controller
package command

import "github.com/mirror-media/major-tom-go/config"

// FIXME we need proper test path
var clusterConfigs = config.K8S{
	"mm": {
		"prod":    "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"staging": "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"dev":     "/Users/chiu/dev/mtv/major-tom-go/configs/config",
	},
	"tv": {
		"prod":    "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"staging": "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"dev":     "/Users/chiu/dev/mtv/major-tom-go/configs/config",
	},
	"readr": {
		"prod": "/Users/chiu/dev/mtv/major-tom-go/configs/config",
		"dev":  "/Users/chiu/dev/mtv/major-tom-go/configs/config",
	},
}
