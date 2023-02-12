// Package presence is the package that contains all of the presence related functions and types.
package presence

import (
	"context"
	"time"

	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/gertd/go-pluralize"
	"golang.org/x/text/message"
)

// NewPresence creates a new presence object.
func NewPresence(l logger.Logger, mp *message.Printer, pl *pluralize.Client) *Presence {
	return &Presence{
		l:  l,
		mp: mp,
		pl: pl,

		t: time.NewTicker(time.Hour),
	}
}

// Presence is an object that is used to update the bot's presence.
type Presence struct {
	// l is the logger.
	l logger.Logger
	// mp is the message printer that is used to format the presence.
	mp *message.Printer
	// pl is the pluralize client that is used to pluralize the presence.
	pl *pluralize.Client

	// t is the ticker that is used to update the presence.
	t *time.Ticker
}

// Update updates the presence on schedule.
func (p *Presence) Update(ctx context.Context, s *state.State) {
	for {
		p.update(ctx, s)

		select {
		case <-p.t.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

// update updates the presence.
func (p *Presence) update(ctx context.Context, s *state.State) {
	g, err := s.Guilds()
	if err != nil {
		p.l.Error("failed to get guilds", "error", err)

		return
	}

	gs := len(g)

	if err := s.Gateway().Send(ctx, &gateway.UpdatePresenceCommand{
		Activities: []discord.Activity{{
			Name: p.mp.Sprintf("for %d %s", gs, p.pl.Pluralize("server", gs, false)),
			Type: discord.WatchingActivity,
		}},
	}); err != nil {
		p.l.Error("failed to update presence", "error", err)

		return
	}

	p.l.Info("updated presence", "guilds", gs)
}
