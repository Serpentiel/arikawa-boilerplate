// Package cmd is the package that contains all of the command handling logic.
package cmd

import (
	"context"

	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/builder"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/state"
)

// ping is the command that replies with a pong.
var ping = &Command{
	CreateCommandData: api.CreateCommandData{
		Name:        "ping",
		Description: "Replies with a pong",
	},
	HandlerFunc: func(cmd *Command, s *state.State) cmdroute.CommandHandlerFunc {
		return func(ctx context.Context, _ cmdroute.CommandData) *api.InteractionResponseData {
			return builder.NewMessageResponse(ctx, cmd.l, cmd.cc, cmd.hc).
				Embed(cmd.mp.Sprintf(
					"🏓 Pong! Bot's latency to Discord is %dms.", s.Gateway().Latency().Milliseconds(),
				)).
				Build()
		}
	},
}
