// Package cmd is the package that contains all of the command handling logic.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/builder"
	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/util"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

// echo is the command that replies with the given text.
var echo = &Command{
	CreateCommandData: api.CreateCommandData{
		Name:        "echo",
		Description: "Replies with the given text",
		Options: discord.CommandOptions{
			&discord.SubcommandOption{
				OptionName:  "normal",
				Description: "Replies with the given text",
				Options: []discord.CommandOptionValue{
					&discord.StringOption{
						OptionName:   "text",
						Description:  "The text to echo",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
			&discord.SubcommandOption{
				OptionName:  "reverse",
				Description: "Reverses the given text and replies with it",
				Options: []discord.CommandOptionValue{
					&discord.StringOption{
						OptionName:   "text",
						Description:  "The text to reverse and echo",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},
	},
	Subs: map[string]*Command{
		"normal": {
			HandlerFunc: func(cmd *Command, s *state.State) cmdroute.CommandHandlerFunc {
				return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
					return builder.NewMessageResponse(ctx, cmd.l, cmd.cc, cmd.hc).Embed(
						data.Options.Find("text").String(),
					).Build()
				}
			},
			AutocompleterFunc: echoAutocompleterFunc,
		},
		"reverse": {
			HandlerFunc: func(cmd *Command, s *state.State) cmdroute.CommandHandlerFunc {
				return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
					return builder.NewMessageResponse(ctx, cmd.l, cmd.cc, cmd.hc).Embed(
						util.StringReverse(data.Options.Find("text").String()),
					).Build()
				}
			},
			AutocompleterFunc: echoAutocompleterFunc,
		},
	},
}

var echoAutocompleterFunc = func(cmd *Command, s *state.State) cmdroute.AutocompleterFunc {
	return func(ctx context.Context, data cmdroute.AutocompleteData) api.AutocompleteChoices {
		const maxAutocompleteChoices = 25

		if data.Options.Focused().Name != "text" {
			return nil
		}

		var choices api.AutocompleteStringChoices

		opt := data.Options.Find("text").String()
		if len(opt) == 0 {
			return choices
		}

		cacheKey := "cmd.echo.autocomplete." + opt

		if v, err := cmd.cc.Any.Get(ctx, cacheKey); err == nil {
			choices, _ = v.(api.AutocompleteStringChoices)

			return choices
		}

		type suggestion struct {
			Word string `json:"word"`
		}

		resp, err := cmd.hc.Get(fmt.Sprintf("https://api.datamuse.com/sug?s=%s", url.QueryEscape(opt)))
		if err != nil {
			cmd.l.Error("failed to get suggestions", "error", err)

			return nil
		}
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				cmd.l.Error("failed to close response body", "error", err)
			}
		}()

		suggestions := []suggestion{}

		if err = json.NewDecoder(resp.Body).Decode(&suggestions); err != nil {
			cmd.l.Error("failed to decode suggestions", "error", err)

			return nil
		}

		choices = make(
			api.AutocompleteStringChoices, util.OrderedMin(len(suggestions), maxAutocompleteChoices),
		)

		for idx, suggestion := range suggestions {
			choices[idx] = discord.StringChoice{
				Name:  suggestion.Word,
				Value: suggestion.Word,
			}
		}

		if err := cmd.cc.Any.Set(ctx, cacheKey, choices); err != nil {
			cmd.l.Error("failed to set cache value", "error", err)
		}

		return choices
	}
}
