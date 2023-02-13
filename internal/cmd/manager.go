// Package cmd is the package that contains all of the command handling logic.
package cmd

import (
	"net/http"

	"github.com/Serpentiel/arikawa-boilerplate/internal/container"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"golang.org/x/text/message"
)

// NewManager creates a new command manager.
func NewManager(l logger.Logger, cc *container.Cache, hc *http.Client, mp *message.Printer) *Manager {
	cm := &Manager{
		l:  l,
		cc: cc,
		hc: hc,
		mp: mp,

		cmds: make(map[string]*Command),
	}

	cm.RegisterCommands(
		ping,
		echo,
	)

	return cm
}

// Manager is the manager for commands.
type Manager struct {
	// l is the logger instance.
	l logger.Logger
	// cc is the container.Cache instance.
	cc *container.Cache
	// hc is the http.Client instance.
	hc *http.Client
	// mp is the message.Printer instance.
	mp *message.Printer

	// cmds is the map of commands.
	cmds map[string]*Command
}

// All returns all of the commands.
func (m *Manager) All() map[string]*Command {
	return m.cmds
}

// injectDependencies injects the dependencies into the command.
func (m *Manager) injectDependencies(cmd *Command) {
	cmd.l = m.l
	cmd.cc = m.cc
	cmd.hc = m.hc
	cmd.mp = m.mp

	for _, sub := range cmd.Subs {
		m.injectDependencies(sub)
	}
}

// RegisterCommand registers a command.
func (m *Manager) RegisterCommand(cmd *Command) {
	m.injectDependencies(cmd)

	m.cmds[cmd.Name] = cmd
}

// RegisterCommands registers multiple commands.
func (m *Manager) RegisterCommands(cmds ...*Command) {
	for _, cmd := range cmds {
		m.RegisterCommand(cmd)
	}
}
