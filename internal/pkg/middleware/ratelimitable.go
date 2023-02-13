// Package middleware provides a set of middleware for the router.
package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/builder"
	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/container"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"golang.org/x/time/rate"
)

// RateLimitable returns a middleware that rate limits users.
func RateLimitable(l logger.Logger, cc *container.Cache, hc *http.Client) cmdroute.Middleware {
	const (
		// rateLimitSeconds is the amount of seconds a user is rate limited for.
		rateLimitSeconds = 3

		// rateLimitBurst is the amount of command interactions a user can send before being rate limited.
		rateLimitBurst = 1
	)

	userLimiters := &sync.Map{}

	return func(next cmdroute.InteractionHandler) cmdroute.InteractionHandler {
		return cmdroute.InteractionHandlerFunc(
			func(ctx context.Context, e *discord.InteractionEvent) *api.InteractionResponse {
				if e.Data.InteractionType() != discord.CommandInteractionType {
					return next.HandleInteraction(ctx, e)
				}

				v, ok := userLimiters.Load(e.Member.User.ID)
				if !ok {
					v, _ = userLimiters.LoadOrStore(
						e.Member.User.ID,
						rate.NewLimiter(rate.Every(rateLimitSeconds*time.Second), rateLimitBurst),
					)
				}

				limiter, _ := v.(*rate.Limiter)

				if !limiter.Allow() {
					return &api.InteractionResponse{
						Type: api.MessageInteractionWithSource,
						Data: builder.NewMessageResponse(ctx, l, cc, hc).
							Error("You are being rate limited! Please wait a few seconds before trying again.").
							Build(),
					}
				}

				return next.HandleInteraction(ctx, e)
			},
		)
	}
}
