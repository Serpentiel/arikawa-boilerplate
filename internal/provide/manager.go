// Package provide is the package that contains all of the providers for the dependency injection container.
package provide

import (
	"context"
	"net/http"

	"github.com/Serpentiel/arikawa-boilerplate/internal/cmd"
	"github.com/Serpentiel/arikawa-boilerplate/internal/container"
	"github.com/Serpentiel/arikawa-boilerplate/internal/presence"
	"github.com/Serpentiel/arikawa-boilerplate/internal/router"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"golang.org/x/text/message"
)

// Manager is a function which provides a *shard.Manager instance.
func Manager(
	lc fx.Lifecycle,
	v *viper.Viper,
	l logger.Logger,
	cc *container.Cache,
	hc *http.Client,
	mp *message.Printer,
	p *presence.Presence,
	cm *cmd.Manager,
) (*shard.Manager, error) {
	isFirstShard := true

	m, err := shard.NewManager(
		"Bot "+v.GetString("discord.bot_token"),
		state.NewShardFunc(func(m *shard.Manager, s *state.State) {
			s.AddIntents(gateway.IntentGuilds)

			s.AddHandler(func(e *gateway.ReadyEvent) {
				if isFirstShard {
					u, err := s.Me()
					if err != nil {
						l.Fatal("failed to get self user", "error", err)
					}

					l.Info("got self user", "tag", u.Tag())

					isFirstShard = false
				}

				l.Info("shard is ready", "id", e.Shard.ShardID(), "num", e.Shard.NumShards())

				go p.Update(s.Context(), s)
			})

			r, err := router.NewRouter(l, cc, hc, mp, cm, s)
			if err != nil {
				l.Fatal("failed to create router", "error", err)
			}

			s.AddInteractionHandler(r)
		}),
	)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return m.Open(ctx)
		},
		OnStop: func(context.Context) error {
			return m.Close()
		},
	})

	return m, nil
}
