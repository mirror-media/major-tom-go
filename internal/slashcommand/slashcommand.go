// Package slashcommand bridge the slack and operations in different projects
package slashcommand

import (
	"context"
	"strings"

	"github.com/mirror-media/major-tom-go/v2/config"
	"github.com/mirror-media/major-tom-go/v2/internal/command"

	"github.com/pkg/errors"
)

const ACCEPTED_SLASHCMD = "/major-tom"

// Response contains instruction of the slashcommand api operation
type CMD struct {
	Command string
	Text    string
}

// Run perform operation per cmd and txt. ctx is expected to have a response channel
func Run(ctx context.Context, clusterConfigs config.K8S, cmd string, txt string) (messages []string, err error) {
	if cmd != ACCEPTED_SLASHCMD {
		return []string{"call help"}, errors.Errorf("%s is not a supported slash command", cmd)
	}
	txtParts := strings.Split(txt, " ")
	if len(txtParts) == 0 {
		// TODO send help
		return []string{"call help"}, nil
	}
	switch cmd := txtParts[0]; cmd {
	case "list":
		messages, err := command.List(ctx, clusterConfigs, txtParts[1:])
		return messages, err
	}
	// TODO send help
	return []string{"call help"}, errors.Errorf("command(%s) is not supported", cmd)
}
