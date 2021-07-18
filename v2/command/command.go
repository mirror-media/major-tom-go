// Package command implement the operation controller
package command

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// response is supposed to use in a channel to pass response for different go routine
type response struct {
	Messages []string
	Error    error
}

func pop(slice []string, i int) (string, []string) {
	ret := slice[i]
	return ret, append(slice[:i], slice[i+1:]...)
}

func cmpArg(text, arg, delimeter string) bool {
	return arg == strings.Split(text, delimeter)[0]
}

func popValue(textParts []string, arg, delimeter string) (newTextParts []string, value string, err error) {
	var result string
	for i, pair := range textParts {
		if cmpArg(pair, arg, delimeter) {
			result, textParts = pop(textParts, i)
			break
		}
	}
	if result == "" {
		return textParts, "", errors.New(fmt.Sprintf("argument(%s) is expected", arg))
	}
	value = strings.Split(result, delimeter)[1]

	return textParts, value, nil
}
