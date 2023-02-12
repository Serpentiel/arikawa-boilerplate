// Package cmd is the package that contains all of the command handling logic.
package cmd

import (
	"net/http"

	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/container"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/state"
	"golang.org/x/text/message"
)

// Command is the struct that contains all of the information about a command.
type Command struct {
	// l is the logger instance.
	l logger.Logger
	// cc is the container.Cache instance.
	cc *container.Cache
	// hc is the http.Client instance.
	hc *http.Client
	// mp is the message.Printer instance.
	mp *message.Printer

	api.CreateCommandData

	// Subs is the map of subcommands.
	Subs map[string]*Command

	// HandlerFunc is the command handler function.
	HandlerFunc func(cmd *Command, s *state.State) cmdroute.CommandHandlerFunc
	// AutocompleterFunc is the command autocompleter function.
	AutocompleterFunc func(cmd *Command, s *state.State) cmdroute.AutocompleterFunc
}
