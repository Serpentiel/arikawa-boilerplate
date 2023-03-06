// Package builder is the package that contains all of the builder functions and types.
package builder

import (
	"context"
	"net/http"

	"github.com/Serpentiel/arikawa-boilerplate/internal/container"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

// NewMessageResponse returns a new MessageResponse.
func NewMessageResponse(
	ctx context.Context,
	l logger.Logger,
	cc *container.Cache,
	hc *http.Client,
) *MessageResponse {
	return &MessageResponse{
		ctx: ctx,
		l:   l,
		cc:  cc,
		hc:  hc,

		d: &api.InteractionResponseData{},
	}
}

// MessageResponse is a struct that contains logic for responding to messages.
type MessageResponse struct {
	// ctx is the context.Context.
	ctx context.Context
	// l is the logger.
	l logger.Logger
	// cc is the container.Cache.
	cc *container.Cache
	// hc is the http.Client.
	hc *http.Client

	// d is the InteractionResponseData.
	d *api.InteractionResponseData
}

// Build returns the InteractionResponseData.
func (b *MessageResponse) Build() *api.InteractionResponseData {
	return b.d
}

// Embed adds an embed to the response.
func (b *MessageResponse) Embed(msg string) *MessageResponse {
	if b.d.Embeds == nil {
		b.d.Embeds = &[]discord.Embed{}
	}

	*b.d.Embeds = append(
		*b.d.Embeds,
		NewEmbed(b.ctx, b.l, b.cc, b.hc).
			Description(msg).
			WithRandomKaomojiTitle().
			WithRandomColor().
			WrapDescriptionInCodeBlock().
			WithRandomWaifuPicture().
			Build(),
	)

	return b
}

// Ephemeral makes the response ephemeral.
func (b *MessageResponse) Ephemeral() *MessageResponse {
	b.d.Flags |= discord.EphemeralMessage

	return b
}

// NoMentions makes the response not mention anyone.
func (b *MessageResponse) NoMentions() *MessageResponse {
	b.d.AllowedMentions = &api.AllowedMentions{}

	return b
}

// Error returns an error embed. It is ephemeral and does not mention anyone.
func (b *MessageResponse) Error(msg string) *MessageResponse {
	return b.Embed(msg).Ephemeral().NoMentions()
}
