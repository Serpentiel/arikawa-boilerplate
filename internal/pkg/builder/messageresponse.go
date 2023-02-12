// Package builder is the package that contains all of the builder functions and types.
package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Serpentiel/arikawa-boilerplate/internal/pkg/cachecontainer"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

// kaomojis is a list of kaomojis.
var kaomojis = []string{
	"(\\* ^ ω ^)",
	"(´ ∀ \\` \\*)",
	"٩(◕‿◕｡)۶",
	"☆\\*:.｡.o(≧▽≦)o.｡.:\\*☆",
	"(o^▽^o)",
	"(⌒▽⌒)☆",
	"<(￣︶￣)>",
	"。.:☆\\*:･'(\\*⌒―⌒\\*)))",
	"ヽ(・∀・)ﾉ",
	"(´｡• ω •｡\\`)",
	"(￣ω￣)",
	"｀;:゛;｀;･(°ε° )",
	"(o･ω･o)",
	"(＠＾◡＾)",
	"ヽ(\\*・ω・)ﾉ",
	"(o\\_ \\_)ﾉ彡☆",
	"(^人^)",
	"(o´▽\\`o)",
	"(\\*´▽\\`\\*)",
	"｡ﾟ( ﾟ^∀^ﾟ)ﾟ｡",
	"( ´ ω \\` )",
	"(((o(\\*°▽°\\*)o)))",
	"(≧◡≦)",
	"(o´∀\\`o)",
	"(´• ω •\\`)",
	"(＾▽＾)",
	"(⌒ω⌒)",
	"∑d(°∀°d)",
	"╰(▔∀▔)╯",
	"(─‿‿─)",
	"(\\*^‿^\\*)",
	"ヽ(o^ ^o)ﾉ",
	"(✯◡✯)",
	"(◕‿◕)",
	"(\\*≧ω≦\\*)",
	"(☆▽☆)",
	"(⌒‿⌒)",
	"＼(≧▽≦)／",
	"ヽ(o＾▽＾o)ノ",
	"☆ ～('▽^人)",
	"(\\*°▽°\\*)",
	"٩(｡•́‿•̀｡)۶",
	"(✧ω✧)",
	"ヽ(\\*⌒▽⌒\\*)ﾉ",
	"(´｡• ᵕ •｡\\`)",
	"( ´ ▽ \\` )",
	"(￣▽￣)",
	"╰(\\*´︶\\`\\*)╯",
	"ヽ(>∀<☆)ノ",
	"o(≧▽≦)o",
	"(☆ω☆)",
	"(っ˘ω˘ς )",
	"＼(￣▽￣)／",
	"(\\*¯︶¯\\*)",
	"＼(＾▽＾)／",
	"٩(◕‿◕)۶",
	"(o˘◡˘o)",
	"\\(★ω★)/",
	"\\(^ヮ^)/",
	"(〃＾▽＾〃)",
	"(╯✧▽✧)╯",
	"o(>ω<)o",
	"o( ❛ᴗ❛ )o",
	"｡ﾟ(TヮT)ﾟ｡",
	"( ‾́ ◡ ‾́ )",
	"(ﾉ´ヮ\\`)ﾉ\\*: ･ﾟ",
	"(b ᵔ▽ᵔ)b",
	"(๑˃ᴗ˂)ﻭ",
	"(๑˘︶˘๑)",
	"( ˙꒳˙ )",
	"(\\*꒦ິ꒳꒦ີ)",
	"°˖✧◝(⁰▿⁰)◜✧˖°",
	"(´･ᴗ･ \\` )",
	"(ﾉ◕ヮ◕)ﾉ\\*:･ﾟ✧",
	"(„• ֊ •„)",
	"(.❛ ᴗ ❛.)",
	"(⁀ᗢ⁀)",
	"(￢‿￢ )",
	"(¬‿¬ )",
	"(\\*￣▽￣)b",
	"( ˙▿˙ )",
	"(¯▿¯)",
	"( ◕▿◕ )",
	"＼(٥⁀▽⁀ )／",
	"(„• ᴗ •„)",
	"(ᵔ◡ᵔ)",
	"( ´ ▿ \\` )",
}

// colors is a slice of uint32s that represent the colors of the embeds.
var colors = []discord.Color{
	12644496, // #c0f090
	11047152, // #a890f0
	16754856, // #ffa8a8
	16773264, // #fff090
	9486591,  // #90c0ff
	15771816, // #f0a8a8
	15790224, // #f0f090
	9486576,  // #90c0f0
}

// NewMessageResponse returns a new MessageResponse.
func NewMessageResponse(
	ctx context.Context,
	l logger.Logger,
	cc *cachecontainer.CacheContainer,
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
	// cc is the cachecontainer.CacheContainer.
	cc *cachecontainer.CacheContainer
	// hc is the http.Client.
	hc *http.Client

	// d is the InteractionResponseData.
	d *api.InteractionResponseData
}

// Build returns the InteractionResponseData.
func (r *MessageResponse) Build() *api.InteractionResponseData {
	return r.d
}

// waifuPicture returns a URL to a waifu picture.
func (r *MessageResponse) waifuPicture() (url string, err error) {
	resp, err := r.hc.Get("https://api.waifu.pics/sfw/waifu")
	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	type waifu struct {
		URL string `json:"url"`
	}

	w := &waifu{}

	if err = json.NewDecoder(resp.Body).Decode(w); err != nil {
		return
	}

	return w.URL, nil
}

// Embed adds an embed to the response.
// nolint:gosec // rand.Intn is not used for security purposes.
func (r *MessageResponse) Embed(msg string) *MessageResponse {
	if r.d.Embeds == nil {
		r.d.Embeds = &[]discord.Embed{}
	}

	e := discord.Embed{
		Title:       kaomojis[rand.Intn(len(kaomojis))],
		Description: fmt.Sprintf("```\n%s\n```", msg),
		Color:       colors[rand.Intn(len(colors))],
	}

	pic, err := r.waifuPicture()
	if err != nil {
		r.l.Error("failed to get waifu picture, retrieving from cache", "error", err)

		pic, err = r.cc.String.Get(r.ctx, "waifu")
		if err != nil {
			r.l.Error("failed to get cache value", "error", err)
		}
	}

	if pic != "" {
		e.Thumbnail = &discord.EmbedThumbnail{URL: pic}

		if err := r.cc.String.Set(r.ctx, "waifu", pic); err != nil {
			r.l.Error("failed to set cache value", "error", err)
		}
	}

	*r.d.Embeds = append(*r.d.Embeds, e)

	return r
}

// Ephemeral makes the response ephemeral.
func (r *MessageResponse) Ephemeral() *MessageResponse {
	r.d.Flags |= discord.EphemeralMessage

	return r
}

// NoMentions makes the response not mention anyone.
func (r *MessageResponse) NoMentions() *MessageResponse {
	r.d.AllowedMentions = &api.AllowedMentions{}

	return r
}

// Error returns an error embed. It is ephemeral and does not mention anyone.
func (r *MessageResponse) Error(msg string) *MessageResponse {
	return r.Embed(msg).Ephemeral().NoMentions()
}
