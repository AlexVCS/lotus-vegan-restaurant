// Package site is the Lotus Vegan website rendered through GoFastr's UI
// system: a core-ui app with one site layout (the owned siteheader and
// sitefooter packages around the primary slot) and three screens, every
// page composed from framework/ui components and themed by
// internal/theme.
package site

import (
	"context"
	"fmt"
	"time"

	uiapp "github.com/DonaldMurillo/gofastr/core-ui/app"
	"github.com/DonaldMurillo/gofastr/core-ui/style"
	"github.com/DonaldMurillo/gofastr/core/render"
	ui "github.com/DonaldMurillo/gofastr/framework/ui"
	"github.com/DonaldMurillo/gofastr/framework/uihost"

	"github.com/AlexVCS/lotus-vegan/internal/sitefooter"
	"github.com/AlexVCS/lotus-vegan/internal/siteheader"
	lotustheme "github.com/AlexVCS/lotus-vegan/internal/theme"
)

// The page bands, registered once as scoped theme overrides; ui.Themed
// wraps a band's content in the override's class.
var (
	limeBand   = style.RegisterThemeOverride(lotustheme.LimeBand)
	goldBand   = style.RegisterThemeOverride(lotustheme.GoldBand)
	forestBand = style.RegisterThemeOverride(lotustheme.ForestBand)
)

const (
	siteName        = "Lotus Vegan"
	siteDescription = "Lotus Vegan Restaurant serves peaceful, plant-based Vietnamese and Asian comfort food in Pinellas Park, Florida."
	// assetVersion busts caches for the favicons and the site script.
	assetVersion = "lv-20261018-1"
)

// New declares the site: theme, layout and screens.
func New() *uiapp.App {
	site := uiapp.NewApp("Lotus Vegan Restaurant").
		WithLang("en").
		WithSkipLabel("Skip to content").
		WithTheme(lotustheme.Lotus().Extend(siteheader.Tokens))

	shell := uiapp.NewLayout("site", uiapp.LayoutSpec{}, buildShell)
	site.SetDefaultLayout(shell)

	site.RegisterScreen(uiapp.NewScreen("/", &HomeScreen{}).
		WithTitle("Plant-based Vietnamese comfort food in Pinellas Park").
		WithDescription(siteDescription).
		WithStyle(HomeStyle), nil)
	site.RegisterScreen(uiapp.NewScreen("/menu", &MenuScreen{}).
		WithTitle("Menu").
		WithDescription("Explore the full Lotus Vegan Restaurant menu, including Vietnamese and Asian plant-based dishes, fresh juices, smoothies, and desserts."), nil)
	site.RegisterScreen(uiapp.NewScreen("/contact", &ContactScreen{}).
		WithTitle("Contact").
		WithDescription("Contact Lotus Vegan Restaurant in Pinellas Park, Florida with questions about the menu, dietary needs, or your next visit.").
		WithStyle(ContactStyle), nil)
	return site
}

// HostOptions are the uihost options: head metadata, favicons, the
// self-hosted DM Sans faces and the site script.
func HostOptions() []uihost.Option {
	return []uihost.Option{
		uihost.WithLang("en"),
		uihost.WithThemeColor("#f5f0e5"),
		uihost.WithOpenGraph(uihost.OG{
			Title:       "Lotus Vegan Restaurant",
			Description: "Be vegan, make peace. Vietnamese and Asian plant-based comfort food in Pinellas Park.",
			Image:       "/assets/avocado-salad.jpg",
			Type:        "restaurant",
		}),
		uihost.WithHeadHTML(fmt.Sprintf(
			`<link rel="icon" href="/favicon.ico?v=%[1]s" sizes="16x16 32x32 48x48">`+
				`<link rel="icon" href="/favicon-32x32.png?v=%[1]s" type="image/png" sizes="32x32">`+
				`<link rel="icon" href="/favicon.svg?v=%[1]s" type="image/svg+xml">`+
				`<link rel="apple-touch-icon" href="/apple-touch-icon.png?v=%[1]s" sizes="180x180">`, assetVersion)),
		// DM Sans, self-hosted as two static weights (public/assets).
		// FontFaceCSS keeps one face per family per call, so each weight
		// is its own call; the theme's Fonts tokens name the family.
		uihost.WithCustomCSS(
			style.FontFaceCSS("/assets", style.WebFont{Family: "DM Sans", File: "dm-sans-regular", Weight: "400"}) +
				style.FontFaceCSS("/assets", style.WebFont{Family: "DM Sans", File: "dm-sans-medium", Weight: "500 700"})),
		// lotus.js carries an old "lotus-theme" choice over to
		// gofastr.colorScheme and keeps <meta name="theme-color"> in step
		// with the colour scheme.
		uihost.WithExtraScripts("/lotus.js?v=" + assetVersion),
	}
}

// buildShell is the site layout: the sticky header, the page, and the
// forest-band footer in a page-tall stack.
func buildShell(ctx context.Context, l *uiapp.LayoutTree) render.HTML {
	return ui.Stack(ui.StackConfig{Screen: true, StickyHeader: true, Gap: ui.GapNone},
		siteHeader(ctx),
		l.Primary(),
		ui.Themed(forestBand, siteFooter()),
		ui.BackToTop(ui.BackToTopConfig{}),
	)
}

func siteHeader(ctx context.Context) render.HTML {
	order := make([]siteheader.OrderLink, 0, len(deliveryServices))
	for _, s := range deliveryServices {
		order = append(order, siteheader.OrderLink{Label: "Order on " + s.Name, Href: s.Href, Icon: s.Icon})
	}
	return siteheader.Render(siteheader.Config{
		Ctx:  ctx,
		Name: siteName,
		Mark: render.Raw(lotusMark),
		Links: []siteheader.Link{
			{Label: "Our food", Href: "/#story"},
			{Label: "Menu", Href: "/menu"},
			{Label: "Visit", Href: "/#visit"},
			{Label: "Contact", Href: "/contact"},
		},
		Order:      order,
		OrderLabel: "Order delivery",
		CTA:        siteheader.Link{Label: "Call to order", Href: phones[0].Tel},
		Actions:    ui.ThemeToggle(ui.ThemeToggleConfig{Variant: ui.ThemeToggleIcon, Ctx: ctx}),
	})
}

func siteFooter() render.HTML {
	order := make([]sitefooter.Link, 0, len(deliveryServices))
	for _, s := range deliveryServices {
		order = append(order, sitefooter.Link{Label: s.Name, Href: s.Href, External: true})
	}
	return sitefooter.Render(sitefooter.Config{
		Name:    siteName,
		Mark:    render.Raw(lotusMark),
		Tagline: "Nhà Hàng Chay Hoa Sen · 6575 Park Blvd N, Pinellas Park, FL 33781",
		Columns: []sitefooter.Column{
			{Title: "Order for delivery", Links: order},
			{Title: "Navigate", Links: []sitefooter.Link{
				{Label: "Our food", Href: "/#story"},
				{Label: "Menu", Href: "/menu"},
				{Label: "Visit", Href: "/#visit"},
				{Label: "Contact", Href: "/contact"},
			}},
			{Title: "Follow", Links: []sitefooter.Link{
				{Label: "Instagram", Href: instagramURL, External: true},
				{Label: "Facebook", Href: facebookURL, External: true},
			}},
		},
		Note: fmt.Sprintf("© %d Lotus Vegan Restaurant · Be vegan, make peace.", time.Now().Year()),
	})
}
