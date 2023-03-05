// Package cmd is the package that contains all of the command handling logic.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/builder"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

// animeData is the data structure for the anime data.
type animeData struct {
	// ID is the ID of the anime.
	ID int `json:"mal_id"`
	// URL is the URL of the anime.
	URL string `json:"url"`
	// Images is the images of the anime.
	Images struct {
		// JPG is the JPG images of the anime.
		JPG struct {
			// LargeImageURL is the large image URL of the anime.
			LargeImageURL string `json:"large_image_url"`
		} `json:"jpg"`
	} `json:"images"`
	// Trailer is the trailer of the anime.
	Trailer struct {
		// URL is the URL of the trailer.
		URL string `json:"url"`
	} `json:"trailer"`
	// Titles is the titles of the anime.
	Titles []struct {
		// Type is the type of the title.
		Type string `json:"type"`
		// Title is the title.
		Title string `json:"title"`
	} `json:"titles"`
	// Type is the type of the anime.
	Type string `json:"type"`
	// Score is the score of the anime.
	Score float64 `json:"score"`
	// ScoredBy is the number of people who scored the anime.
	ScoredBy int `json:"scored_by"`
	// Rank is the rank of the anime.
	Rank int `json:"rank"`
	// Popularity is the popularity of the anime.
	Popularity int `json:"popularity"`
	// Members is the number of members who have the anime on their list.
	Members int `json:"members"`
	// Synopsis is the synopsis of the anime.
	Synopsis string `json:"synopsis"`
	// Year is the year of the anime.
	Year int `json:"year"`
	// Studios is the studios of the anime.
	Studios []struct {
		// Name is the name of the studio.
		Name string `json:"name"`
		// URL is the URL of the studio.
		URL string `json:"url"`
	} `json:"studios"`
}

// anime is the command that lets you search for anime.
var anime = &Command{
	CreateCommandData: api.CreateCommandData{
		Name:        "anime",
		Description: "Searches for anime",
		Options: discord.CommandOptions{
			&discord.StringOption{
				OptionName:   "query",
				Description:  "The query to search for",
				Required:     true,
				Autocomplete: true,
			},
		},
	},
	HandlerFunc: func(cmd *Command, s *state.State) cmdroute.CommandHandlerFunc {
		return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
			opt := data.Options.Find("query").String()
			if len(opt) == 0 {
				return builder.NewMessageResponse(ctx, cmd.l, cmd.cc, cmd.hc).
					Error("You must provide a MAL ID to search for.").
					Build()
			}

			var ird *api.InteractionResponseData

			cacheKey := "cmd.anime.handle." + opt

			if v, err := cmd.cc.Any.Get(ctx, cacheKey); err == nil {
				ird, _ = v.(*api.InteractionResponseData)

				return ird
			}

			// result is the result of the anime lookup by ID.
			type result struct {
				// Data is the data of the anime.
				Data animeData `json:"data"`
			}

			resp, err := cmd.hc.Get(fmt.Sprintf("https://api.jikan.moe/v4/anime/%s", url.QueryEscape(opt)))
			if err != nil {
				cmd.l.Error("failed to get anime", "error", err)

				return nil
			}
			defer func() {
				err = resp.Body.Close()
				if err != nil {
					cmd.l.Error("failed to close response body", "error", err)
				}
			}()

			res := result{}

			if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
				cmd.l.Error("failed to decode anime response", "error", err)

				return nil
			}

			d := res.Data

			e := builder.NewEmbed(ctx, cmd.l, cmd.cc, cmd.hc).
				Title(animeTitle(d, false)).
				Description(d.Synopsis).
				URL(d.URL).
				Image(&discord.EmbedImage{URL: d.Images.JPG.LargeImageURL}).
				Fields([]discord.EmbedField{
					{
						Name: "Score",
						Value: cmd.mp.Sprintf(
							"%s (scored by %d)",
							strconv.FormatFloat(d.Score, 'f', -1, 64),
							d.ScoredBy,
						),
						Inline: true,
					},
					{
						Name:   "Rank",
						Value:  cmd.mp.Sprintf("#%d", d.Rank),
						Inline: true,
					},
					{
						Name:   "Popularity",
						Value:  cmd.mp.Sprintf("#%d (watched by %d)", d.Popularity, d.Members),
						Inline: true,
					},
				}).
				WithRandomColor()

			if len(d.Studios) > 0 {
				author := &discord.EmbedAuthor{}

				var names []string

				for _, studio := range d.Studios {
					names = append(names, studio.Name)
				}

				author.Name = strings.Join(names, ", ")

				if len(d.Studios) == 1 {
					author.URL = d.Studios[0].URL
				}

				e = e.Author(author)
			}

			arc := discord.ActionRowComponent{
				&discord.ButtonComponent{
					Label: "View on MAL",
					Style: discord.LinkButtonStyle(d.URL),
				},
			}

			if len(d.Trailer.URL) > 0 {
				arc = append(arc, &discord.ButtonComponent{
					Label: "View Trailer",
					Style: discord.LinkButtonStyle(d.Trailer.URL),
				})
			}

			ird = &api.InteractionResponseData{
				Embeds:     &[]discord.Embed{e.Build()},
				Components: discord.ComponentsPtr(&arc),
			}

			if err := cmd.cc.Any.Set(ctx, cacheKey, ird); err != nil {
				cmd.l.Error("failed to set cache value", "error", err)
			}

			return ird
		}
	},
	AutocompleterFunc: func(cmd *Command, s *state.State) cmdroute.AutocompleterFunc {
		return func(ctx context.Context, data cmdroute.AutocompleteData) api.AutocompleteChoices {
			if data.Options.Focused().Name != "query" {
				return nil
			}

			var choices api.AutocompleteStringChoices

			opt := data.Options.Find("query").String()
			if len(opt) == 0 {
				return choices
			}

			cacheKey := "cmd.anime.autocomplete." + opt

			if v, err := cmd.cc.Any.Get(ctx, cacheKey); err == nil {
				choices, _ = v.(api.AutocompleteStringChoices)

				return choices
			}

			// result is the result of the anime search.
			type result struct {
				// Data is the list of anime data.
				Data []animeData `json:"data"`
			}

			resp, err := cmd.hc.Get(fmt.Sprintf("https://api.jikan.moe/v4/anime?q=%s&limit=25", url.QueryEscape(opt)))
			if err != nil {
				cmd.l.Error("failed to get list of anime", "error", err)

				return nil
			}
			defer func() {
				err = resp.Body.Close()
				if err != nil {
					cmd.l.Error("failed to close response body", "error", err)
				}
			}()

			res := result{}

			if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
				cmd.l.Error("failed to decode list of anime", "error", err)

				return nil
			}

			d := res.Data

			choices = make(api.AutocompleteStringChoices, len(d))

			for idx, anime := range d {
				choices[idx] = discord.StringChoice{
					Name:  animeTitle(anime, true),
					Value: strconv.Itoa(anime.ID),
				}
			}

			if err := cmd.cc.Any.Set(ctx, cacheKey, choices); err != nil {
				cmd.l.Error("failed to set cache value", "error", err)
			}

			return choices
		}
	},
}

// animeTitle returns the title of the anime.
func animeTitle(d animeData, clamp bool) string {
	// titleOffset is the maximum possible length of the title suffixed with "… (type, year)".
	const titleOffset = 17

	title := d.Titles[0].Title
	if clamp && len(title)+titleOffset > maxAutocompleteStringChoiceLength {
		title = title[:maxAutocompleteStringChoiceLength-titleOffset] + "…"
	}

	b := strings.Builder{}

	b.WriteString(title)

	b.WriteString(" (")

	b.WriteString(d.Type)

	if d.Year > 0 {
		b.WriteString(", ")

		b.WriteString(strconv.Itoa(d.Year))
	}

	b.WriteRune(')')

	return b.String()
}
