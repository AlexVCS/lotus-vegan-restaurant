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
	dark := Dark
	o := Light
	o.Dark = &dark
	o.FontBody = FontStack
	o.FontHeading = FontStack

	t := uitheme.Default(o)
	t.Name = "lotus"

	// Secondary buttons sit on the soft cream / deep-green fill.
	t.Colors.Secondary.Value = "#ebe4d6"
	t.Colors.SecondaryFg.Value = "#16372a"
	t.DarkColors["secondary"] = "#1b3329"
	t.DarkColors["secondary-fg"] = "#f7f2e8"

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
