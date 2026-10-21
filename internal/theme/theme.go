// Package theme is the Lotus Vegan brand expressed as a GoFastr theme.
//
// Lotus builds on framework/ui/theme.Default, the adaptive framework
// theme, and swaps in the restaurant's palette: cream pages, forest-green
// ink and buttons, coral emphasis in light mode, and a deep green page
// with lime accents in dark mode. Brand colours the framework palette has
// no slot for (lime, gold, the dietary badges) are app tokens declared in
// brand.tokens.css and generated into brand_tokens.gen.go by
// `gofastr gen styles`.
//
// Every value here is a token: the framework components read it through
// var(--color-…), var(--font-…), var(--shadow-…) and so on, and the dark
// values are re-declared under :root[data-color-scheme="dark"].
package theme

import (
	"github.com/DonaldMurillo/gofastr/core-ui/style"
	uitheme "github.com/DonaldMurillo/gofastr/framework/ui/theme"
)

// FontStack is the body and heading family: the self-hosted DM Sans
// (public/assets/dm-sans-*.ttf) over the platform sans-serif.
const FontStack = `"DM Sans", ui-sans-serif, system-ui, -apple-system, "Segoe UI", sans-serif`

// Light is the Lotus light palette.
var Light = uitheme.Overrides{
	Background:   "#f5f0e5", // cream page
	Surface:      "#fbf7ee", // cards and controls
	SurfaceSoft:  "#ebe4d6", // hover fills, soft chips
	Border:       "#cfd3c8",
	BorderStrong: "#a8b4ad",
	Text:         "#16372a",
	TextMuted:    "#4b6257",
	TextSubtle:   "#5a6e64",
	Primary:      "#173c2d", // forest
	PrimaryFg:    "#fffaf0",
	Accent:       "#b94127", // coral emphasis
	Danger:       "#a6321b",
	DangerFg:     "#ffffff",
	Success:      "#2f5d1e",
	Warning:      "#7a4a06",
	Info:         "#1f5a7a",
}

// Dark is the Lotus dark palette.
var Dark = uitheme.Overrides{
	Background:   "#0a1711",
	Surface:      "#10231b",
	SurfaceSoft:  "#1b3329",
	Border:       "#2c4439",
	BorderStrong: "#52685e",
	Text:         "#f7f2e8",
	TextMuted:    "#c4cec7",
	TextSubtle:   "#a3b3aa",
	Primary:      "#d5ed74", // lime
	PrimaryFg:    "#10231b",
	Accent:       "#ff9b7d",
	Danger:       "#ff9b7d",
	DangerFg:     "#10231b",
	Success:      "#a8d886",
	Warning:      "#f0c46a",
	Info:         "#8cc8e6",
}

// Lotus returns the site theme: the framework's adaptive default with
// the Lotus palette, DM Sans, forest-tinted shadows and the brand
// tokens. It is validated (required tokens and the WCAG AA ink pairs)
// before it is returned, so a bad edit panics at boot rather than
// shipping unreadable buttons.
func Lotus() style.Theme {
	return build(Light, Dark)
}

// The page bands. Each is the full Lotus theme with the band's palette
// swapped in, so a ui.Themed scope re-declares every token (fonts,
// shadows, brand tokens) and the components inside it read the band's
// colours. Every band carries a dark palette, so it follows the
// visitor's scheme like the page does.
var (
	// LimeBand is the quick-facts and visit band: lime in light mode,
	// a quiet deep green in dark mode.
	LimeBand = band(uitheme.Overrides{
		Background: "#cae75d", Surface: "#d6ee7c", SurfaceSoft: "#bcd94f",
		Border: "#a3bd45", BorderStrong: "#7d9433",
		Text: "#173c2d", TextMuted: "#304a20", TextSubtle: "#3a5226",
		Primary: "#173c2d", PrimaryFg: "#f8f3e8",
		Accent: "#7f2b16", Danger: "#7f2b16",
	}, uitheme.Overrides{
		Background: "#153126", Surface: "#1a3a2d", SurfaceSoft: "#214537",
		Border: "#2f5143", BorderStrong: "#52685e",
		Text: "#f7f2e8", TextMuted: "#c8d6cd", TextSubtle: "#aebfb5",
		Primary: "#d5ed74", PrimaryFg: "#10231b",
		Accent: "#ff9b7d", Danger: "#ff9b7d", DangerFg: "#10231b",
	})

	// GoldBand is the pizza band.
	GoldBand = band(uitheme.Overrides{
		Background: "#d8a332", Surface: "#e2b552", SurfaceSoft: "#cf9726",
		Border: "#b8892a", BorderStrong: "#8f6a1f",
		Text: "#173c2d", TextMuted: "#2a3b1f", TextSubtle: "#2f3f22",
		Primary: "#173c2d", PrimaryFg: "#fff7e6",
		Accent: "#173c2d", Danger: "#6e2412",
	}, uitheme.Overrides{
		Background: "#5a4219", Surface: "#674c1f", SurfaceSoft: "#735626",
		Border: "#7d6430", BorderStrong: "#a2854a",
		Text: "#fff7e4", TextMuted: "#f2e3bf", TextSubtle: "#e6d3a5",
		Primary: "#f4d67f", PrimaryFg: "#2b1f08",
		Accent: "#f4d67f", Danger: "#ffb39b", DangerFg: "#2b1f08",
	})

	// ForestBand is the menu band and the footer: forest green with
	// cream ink and lime accents in both schemes.
	ForestBand = band(uitheme.Overrides{
		Background: "#173c2d", Surface: "#1d4636", SurfaceSoft: "#24503f",
		Border: "#355a4a", BorderStrong: "#61776b",
		Text: "#f8f3e8", TextMuted: "#d2dbd4", TextSubtle: "#b9c7be",
		Primary: "#cae75d", PrimaryFg: "#173c2d",
		Accent: "#cae75d", Danger: "#ff9b7d", DangerFg: "#173c2d",
		Success: "#a8d886", Warning: "#f0c46a", Info: "#8cc8e6",
	}, uitheme.Overrides{
		Background: "#0c2a1f", Surface: "#10342a", SurfaceSoft: "#163d31",
		Border: "#294a3d", BorderStrong: "#536a5f",
		Text: "#f7f2e8", TextMuted: "#c8d2cb", TextSubtle: "#a9b8af",
		Primary: "#d5ed74", PrimaryFg: "#0c2a1f",
		Accent: "#d5ed74", Danger: "#ff9b7d", DangerFg: "#0c2a1f",
		Success: "#a8d886", Warning: "#f0c46a", Info: "#8cc8e6",
	})
)

// band builds a band theme: the Lotus palettes with the band's colours
// laid over them.
func band(light, dark uitheme.Overrides) style.Theme {
	l, d := Light, Dark
	mergeColours(&l, light)
	mergeColours(&d, dark)
	return build(l, d)
}

// mergeColours copies every non-empty colour of src over dst.
func mergeColours(dst *uitheme.Overrides, src uitheme.Overrides) {
	for _, f := range []struct {
		d *string
		s string
	}{
		{&dst.Background, src.Background}, {&dst.Surface, src.Surface}, {&dst.SurfaceSoft, src.SurfaceSoft},
		{&dst.Border, src.Border}, {&dst.BorderStrong, src.BorderStrong},
		{&dst.Text, src.Text}, {&dst.TextMuted, src.TextMuted}, {&dst.TextSubtle, src.TextSubtle},
		{&dst.Primary, src.Primary}, {&dst.PrimaryFg, src.PrimaryFg}, {&dst.Accent, src.Accent},
		{&dst.Success, src.Success}, {&dst.Warning, src.Warning}, {&dst.Danger, src.Danger},
		{&dst.DangerFg, src.DangerFg}, {&dst.Info, src.Info},
	} {
		if f.s != "" {
			*f.d = f.s
		}
	}
}

func build(light, dark uitheme.Overrides) style.Theme {
	o := light
	o.Dark = &dark
	o.FontBody = FontStack
	o.FontHeading = FontStack

	t := uitheme.Default(o)
	t.Name = "lotus"

	// Secondary buttons sit on the soft fill of their palette.
	t.Colors.Secondary.Value = light.SurfaceSoft
	t.Colors.SecondaryFg.Value = light.Text
	t.DarkColors["secondary"] = dark.SurfaceSoft
	t.DarkColors["secondary-fg"] = dark.Text

	// Shadows tinted with the forest ink instead of neutral black.
	t.Shadows.XS.Value = "0 1px 2px 0 rgba(13,43,32,0.06)"
	t.Shadows.SM.Value = "0 1px 3px 0 rgba(13,43,32,0.10), 0 1px 2px -1px rgba(13,43,32,0.08)"
	t.Shadows.MD.Value = "0 6px 14px -4px rgba(13,43,32,0.14), 0 2px 4px -2px rgba(13,43,32,0.08)"
	t.Shadows.LG.Value = "0 18px 40px -12px rgba(13,43,32,0.22), 0 4px 8px -4px rgba(13,43,32,0.10)"
	t.Shadows.XL.Value = "0 28px 60px -16px rgba(13,43,32,0.28), 0 8px 12px -6px rgba(13,43,32,0.10)"

	t = t.Extend(Tokens)
	t.MustValidate()
	return t
}
