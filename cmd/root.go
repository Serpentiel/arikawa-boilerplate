// Package cmd is the package that contains all of the commands for the application.
package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"

	dcmd "github.com/Serpentiel/arikawa-boilerplate/internal/cmd"
	"github.com/Serpentiel/arikawa-boilerplate/internal/container"
	"github.com/Serpentiel/arikawa-boilerplate/internal/presence"
	"github.com/Serpentiel/arikawa-boilerplate/internal/provide"
	"github.com/Serpentiel/arikawa-boilerplate/internal/router"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"golang.org/x/text/message"
)

// rootCmd is the rootCmd command for the application.
var rootCmd = &cobra.Command{
	Use: "arikawa-boilerplate",
	Run: func(cmd *cobra.Command, args []string) {
		fx.New(
			fx.Supply(cmd),

			fx.Provide(
				provide.Viper,
				provide.Logger,

				provide.RistrettoStore,
				provide.Cache,

				provide.HTTPClient,

				provide.MessagePrinter,
				pluralize.NewClient,

				presence.NewPresence,
				dcmd.NewManager,

				provide.HTTPServer,
				provide.Manager,
			),

			fx.WithLogger(provide.FxeventLogger),

			fx.Invoke(func(
				lc fx.Lifecycle,
				v *viper.Viper,
				l logger.Logger,
				cc *container.Cache,
				hc *http.Client,
				mp *message.Printer,
				cm *dcmd.Manager,
				hs *http.Server,
				m *shard.Manager,
			) {
				if v.GetBool("http.enabled") {
					withHTTPServer(lc, v, l, cc, hc, mp, cm, hs)
				} else {
					withManager(lc, l, m)
				}
			}),
		).Run()
	},
}

// withHTTPServer is a function that attaches the HTTP server to the lifecycle.
func withHTTPServer(
	lc fx.Lifecycle,
	v *viper.Viper,
	l logger.Logger,
	cc *container.Cache,
	hc *http.Client,
	mp *message.Printer,
	cm *dcmd.Manager,
	hs *http.Server,
) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) (err error) {
			s := state.NewAPIOnlyState("Bot "+v.GetString("discord.bot.token"), nil)

			r, err := router.NewRouter(l, cc, hc, mp, cm, s)
			if err != nil {
				l.Fatal("failed to create router", "error", err)
			}

			is, err := webhook.NewInteractionServer(v.GetString("discord.app.public_key"), r)
			if err != nil {
				l.Fatal("failed to create interaction server", "error", err)
			}

			hs.Handler = is

			go func() {
				if err = hs.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					l.Error("HTTP server error", "error", err)
				}
			}()

			l.Info("HTTP server listening", "address", hs.Addr)

			return
		},
		OnStop: func(ctx context.Context) error {
			return hs.Shutdown(ctx)
		},
	})
}

// withManager is a function that attaches the shard manager to the lifecycle.
func withManager(lc fx.Lifecycle, l logger.Logger, m *shard.Manager) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return m.Open(ctx)
		},
		OnStop: func(context.Context) error {
			return m.Close()
		},
	})
}

// Execute is the function that is called to execute the rootCmd command.
func Execute() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "defines the config file to use")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
