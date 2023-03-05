// Package builder is the package that contains all of the builder functions and types.
package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Serpentiel/arikawa-boilerplate/internal/container"
	"github.com/Serpentiel/arikawa-boilerplate/pkg/logger"
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

// Embed is a struct that contains logic for building embeds.
type Embed struct {
	// ctx is the context.Context.
	ctx context.Context
	// l is the logger.
	l logger.Logger
	// cc is the container.Cache.
	cc *container.Cache
	// hc is the http.Client.
	hc *http.Client

	// e is the discord.Embed.
	e discord.Embed

	// withRandomKaomojiTitle is a bool that determines whether or not to include a random kaomoji in the title.
	withRandomKaomojiTitle bool
	// withRandomColor is a bool that determines whether or not to include a random color.
	withRandomColor bool
	// wrapDescriptionInCodeBlock is a bool that determines whether or not to wrap the description in a code block.
	wrapDescriptionInCodeBlock bool
	// withRandomWaifuPicture is a bool that determines whether or not to include a random waifu picture.
	withRandomWaifuPicture bool
}

// NewEmbed returns a new Embed.
func NewEmbed(ctx context.Context, l logger.Logger, cc *container.Cache, hc *http.Client) *Embed {
	return &Embed{
		ctx: ctx,
		l:   l,
		cc:  cc,
		hc:  hc,

		e: discord.Embed{},
	}
}

// Title sets the title of the embed.
func (b *Embed) Title(title string) *Embed {
	b.e.Title = title

	return b
}

// Description sets the description of the embed.
func (b *Embed) Description(description string) *Embed {
	b.e.Description = description

	return b
}

// URL sets the URL of the embed.
func (b *Embed) URL(url string) *Embed {
	b.e.URL = url

	return b
}

// Color sets the color of the embed.
func (b *Embed) Color(color discord.Color) *Embed {
	b.e.Color = color

	return b
}

// Image sets the image of the embed.
func (b *Embed) Image(image *discord.EmbedImage) *Embed {
	b.e.Image = image

	return b
}

// Author sets the author of the embed.
func (b *Embed) Author(author *discord.EmbedAuthor) *Embed {
	b.e.Author = author

	return b
}

// Fields sets the fields of the embed.
func (b *Embed) Fields(fields []discord.EmbedField) *Embed {
	b.e.Fields = fields

	return b
}

// WithRandomKaomojiTitle sets the withRandomKaomojiTitle bool to true.
func (b *Embed) WithRandomKaomojiTitle() *Embed {
	b.withRandomKaomojiTitle = true

	return b
}

// WithRandomColor sets the withRandomColor bool to true.
func (b *Embed) WithRandomColor() *Embed {
	b.withRandomColor = true

	return b
}

// WrapDescriptionInCodeBlock sets the wrapDescriptionInCodeBlock bool to true.
func (b *Embed) WrapDescriptionInCodeBlock() *Embed {
	b.wrapDescriptionInCodeBlock = true

	return b
}

// WithRandomWaifuPicture sets the withRandomWaifuPicture bool to true.
func (b *Embed) WithRandomWaifuPicture() *Embed {
	b.withRandomWaifuPicture = true

	return b
}

// waifuPicture returns a URL to a waifu picture.
func (b *Embed) waifuPicture() (url string, err error) {
	resp, err := b.hc.Get("https://api.waifu.pics/sfw/waifu")
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

// Build builds the embed.
// nolint:gosec // rand.Intn is not used for security purposes.
func (b *Embed) Build() discord.Embed {
	if b.withRandomKaomojiTitle {
		b.e.Title = kaomojis[rand.Intn(len(kaomojis))]
	}

	if b.wrapDescriptionInCodeBlock {
		b.e.Description = fmt.Sprintf("```\n%s\n```", b.e.Description)
	}

	if b.withRandomColor {
		b.e.Color = colors[rand.Intn(len(colors))]
	}

	if !b.withRandomWaifuPicture {
		return b.e
	}

	pic, err := b.waifuPicture()
	if err != nil {
		b.l.Error("failed to get waifu picture, retrieving from cache", "error", err)

		pic, err = b.cc.String.Get(b.ctx, "waifu")
		if err != nil {
			b.l.Error("failed to get cache value", "error", err)
		}
	}

	if pic != "" {
		b.e.Thumbnail = &discord.EmbedThumbnail{URL: pic}

		if err := b.cc.String.Set(b.ctx, "waifu", pic); err != nil {
			b.l.Error("failed to set cache value", "error", err)
		}
	}

	return b.e
}
