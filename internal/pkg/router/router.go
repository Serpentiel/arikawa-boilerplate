// Package router is the package that provides the command router.
package router

import (
	"net/http"

	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/cmd"
	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/container"
	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/middleware"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/state"
	"golang.org/x/text/message"
)

// NewRouter creates a new Router instance and registers all of the commands.
func NewRouter(
	l logger.Logger,
	cc *container.Cache,
	hc *http.Client,
	mp *message.Printer,
	cm *cmd.Manager,
	s *state.State,
) (*Router, error) {
	r := newRouter(cmdroute.NewRouter(), s)

	r.Use(
		cmdroute.Deferrable(s, cmdroute.DeferOpts{}),
		middleware.RateLimitable(l, cc, hc),
	)

	cmds := cm.All()

	d := make([]api.CreateCommandData, 0, len(cmds))

	for _, cmd := range cmds {
		d = append(d, cmd.CreateCommandData)

		r.addCommand(cmd.Name, cmd)
	}

	return r, cmdroute.OverwriteCommands(s, d)
}

// newRouter creates a new Router instance.
func newRouter(r *cmdroute.Router, s *state.State) *Router {
	return &Router{
		Router: r,

		s: s,
	}
}

// Router is the stateful router.
type Router struct {
	*cmdroute.Router

	// s is the s.State instance.
	s *state.State
}

// addCommand adds a command to the router.
func (r *Router) addCommand(name string, cmd *cmd.Command) {
	subs := cmd.Subs

	if len(subs) > 0 {
		r.Sub(name, func(cr *cmdroute.Router) {
			nr := newRouter(cr, r.s)

			for name, sub := range subs {
				nr.addCommand(name, sub)
			}
		})

		return
	}

	if cmd.HandlerFunc != nil {
		r.AddFunc(name, cmd.HandlerFunc(cmd, r.s))
	}

	if cmd.AutocompleterFunc != nil {
		r.AddAutocompleterFunc(name, cmd.AutocompleterFunc(cmd, r.s))
	}
}
