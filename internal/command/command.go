// Package command implement the operation controller
package command

// FIXME temporary workaround and it belongs to configurations provided to the functionality
var clusters map[string][]string = map[string][]string{
	"mm": {
		"prod",
		"staging",
		"dev",
	},
	"tv": {
		"prod",
		"staging",
		"dev",
	},
	"readr": {
		"prod",
		"dev",
	},
}
