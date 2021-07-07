// Package command implement the operation controller
package command

// response is supposed to use in a channel to pass response for different go routine
type response struct {
	Messages []string
	Error    error
}
