package command

import (
	"context"
	"strings"
)

// List provide infomation about cluster, stages, and services. It should also provide helm message if input is invalid
func List(ctx context.Context, textParts []string) (message []string, err error) {
	switch len(textParts) {
	case 0:
		var projects []string
		for key := range clusters {
			projects = append(projects, key)
		}
		message = []string{
			"The following projects are available: " + strings.Join(projects, ", "),
		}
		return message, nil
	}
	return nil, nil
}
