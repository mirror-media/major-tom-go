// Package slashcommand bridge the slack and operations in different projects
package slashcommand

import (
	"context"
	"strings"

	"github.com/mirror-media/major-tom-go/v2/command"
	"github.com/mirror-media/major-tom-go/v2/config"

	"github.com/pkg/errors"
)

const ACCEPTED_SLASHCMD = "/major-tom"

// Response contains instruction of the slashcommand api operation
type CMD struct {
	Command string
	Text    string
}

// Run perform operation per cmd and txt. ctx is expected to have a response channel
func Run(ctx context.Context /*clusterConfigs config.K8S,*/, k8sRepoConfig config.KubernetesConfigsRepo, slashcmd, txt, caller string) (messages []string, err error) {
	if slashcmd != ACCEPTED_SLASHCMD {
		return []string{"call help"}, errors.Errorf("%s is not a supported slash command", slashcmd)
	}
	txtParts := strings.Split(txt, " ")
	if len(txtParts) == 0 {
		// TODO send help
		return []string{"call help"}, nil
	}

	cmd := txtParts[0]
	switch cmd {
	// case "list":
	// 	messages, err = command.List(ctx, clusterConfigs, txtParts[1:])
	// case "info":
	// 	messages, err = command.Info(ctx, clusterConfigs, txtParts[1:])
	case "deploy":
		messages, err = command.Deploy(ctx, k8sRepoConfig, txtParts[1:], txt, "+"+caller)
	case "release":
		messages, err = command.Release(ctx, k8sRepoConfig, txtParts[1:], txt, "+"+caller)
	default:
		if isBowie(txtParts) {
			messages = command.Bowie()
		} else {
			messages = []string{"call help"}
			err = errors.Errorf("command(%s) is not supported", cmd)
		}
	}
	return messages, err
}

func isBowie(txtParts []string) bool {

	david := map[string]interface{}{
		"tom":      nil,
		"bowie":    nil,
		"space":    nil,
		"oddity":   nil,
		"assemble": nil,
	}

	for _, text := range txtParts {
		if _, ok := david[strings.ToLower(text)]; ok {
			return true
		}
	}
	return false
}
