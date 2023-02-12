// Package middleware provides a set of middleware for the router.
package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/builder"
	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/container"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
)

// RateLimitable returns a middleware that rate limits users.
func RateLimitable(l logger.Logger, cc *container.Cache, hc *http.Client) cmdroute.Middleware {
	const rateLimitSeconds = 3

	rateLimitedUsers := map[discord.UserID]struct{}{}

	return func(next cmdroute.InteractionHandler) cmdroute.InteractionHandler {
		return cmdroute.InteractionHandlerFunc(
			func(ctx context.Context, e *discord.InteractionEvent) *api.InteractionResponse {
				if e.Data.InteractionType() != discord.CommandInteractionType {
					return next.HandleInteraction(ctx, e)
				}

				if _, ok := rateLimitedUsers[e.SenderID()]; ok {
					return &api.InteractionResponse{
						Type: api.MessageInteractionWithSource,
						Data: builder.NewMessageResponse(ctx, l, cc, hc).
							Error("You are being rate limited! Please wait a few seconds before trying again.").
							Build(),
					}
				}

				rateLimitedUsers[e.SenderID()] = struct{}{}

				t := time.NewTicker(rateLimitSeconds * time.Second)

				go func() {
					<-t.C

					t.Stop()

					delete(rateLimitedUsers, e.SenderID())
				}()

				return next.HandleInteraction(ctx, e)
			},
		)
	}
}
