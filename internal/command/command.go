// Package command implement the operation controller
package command

// CommandResponse is supposed to use in a channel to pass response for different go routine
type CommandResponse struct {
	Messages []string
	Error    error
}
