// Package slashcommand bridge the slack and operations in different projects
package slashcommand

import (
	"context"
	"strings"

	"github.com/mirror-media/major-tom-go/v2/internal/command"

	"github.com/pkg/errors"
)

const ACCEPTED_SLASHCMD = "/major-tom"

// Response contains instruction of the slashcommand api operation
type CMD struct {
	Command string
	Text    string
}

// Response contains messages or error for the slashcommand api operation
type Response struct {
	Messages []string
	Err      error
}

// Run perform operation per cmd and txt. ctx is expected to have a response channel
func Run(ctx context.Context, cmd string, txt string) (resp Response) {
	if cmd != ACCEPTED_SLASHCMD {
		return Response{
			// TODO send help
			Messages: []string{"call help"},
			Err:      errors.Errorf("%s is not a supported slash command", cmd),
		}
	}
	txtParts := strings.Split(txt, " ")
	if len(txtParts) == 0 {
		// TODO send help
		return Response{
			Messages: []string{"call help"},
			Err:      errors.Errorf("%s is not a supported slash command", cmd),
		}
	}
	switch txtParts[0] {
	case "list":
		// TODO do operation
	}
}
