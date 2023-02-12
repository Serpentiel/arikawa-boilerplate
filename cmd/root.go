// Package cmd is the package that contains all of the commands for the application.
package cmd

import (
	"os"

	dcmd "github.com/Serpentiel/arikawa-boilerplate/internal/pkg/cmd"
	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/presence"
	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/provide"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
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
				provide.CacheContainer,

				provide.HTTPClient,

				provide.MessagePrinter,
				pluralize.NewClient,

				presence.NewPresence,
				dcmd.NewManager,
				provide.Manager,
			),

			fx.WithLogger(provide.FxeventLogger),

			fx.Invoke(func(*shard.Manager) {}),
		).Run()
	},
}

// Execute is the function that is called to execute the rootCmd command.
func Execute() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "defines the config file to use")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
