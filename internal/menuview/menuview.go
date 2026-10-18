// Package menuview renders the Lotus menu (internal/menu) with GoFastr
// components: a Collapsible per menu group (one open at a time), a Grid of
// outlined Cards for the items (title, price, description, dietary
// badges), StatusBadge + Tooltip for the dietary codes, a jump strip of
// LinkButtons, and a Callout for the allergy note.
//
// The few things the kit does not draw (the scrolling jump strip, the
// price colour, the group number) live in menuview.style.css, an owned
// style whose every value is a theme token.
package menuview

import (
	"fmt"

	"github.com/DonaldMurillo/gofastr/core-ui/html"
	"github.com/DonaldMurillo/gofastr/core/render"
	ui "github.com/DonaldMurillo/gofastr/framework/ui"

	"github.com/AlexVCS/lotus-vegan/internal/menu"
)

// The dietary badges are status variants, so StatusBadge draws them in
// the kit's soft tone: the colour as text on a tint of itself. They sit
// on the forest menu band, where lime and light coral both clear 4.5:1
// (internal/theme tests hold that).
var dietaryVariants = map[string]ui.StatusVariant{
	"gf": ui.RegisterStatusVariant("dietary-gf", ui.StatusVariantCSS{Color: "{colors.brand-lime}"}),
	"sf": ui.RegisterStatusVariant("dietary-sf", ui.StatusVariantCSS{Color: "{colors.dietary-sf}"}),
}

// GroupID is the element id of a menu group, the target of the jump
// strip ("#menu-appetizers").
func GroupID(slug string) string { return "menu-" + slug }

// Config is what a page passes in.
type Config struct {
	Menu *menu.Menu
	// ItemHeadingLevel is the heading level of each item's title.
	ItemHeadingLevel int
}

// Legend renders the dietary key and the allergy-note link.
func Legend() render.HTML {
	items := []render.HTML{}
	for _, d := range menu.DietaryMarks {
		items = append(items, ui.Cluster(ui.ClusterConfig{Gap: ui.GapSM, NoWrap: true},
			ui.StatusBadge(ui.StatusBadgeConfig{Label: d.Label, Variant: dietaryVariants[d.Key], ExtraAttrs: html.Attrs{"aria-hidden": "true"}}),
			html.Span(html.TextConfig{}, render.Text(d.Text)),
		))
	}
	items = append(items, ui.Link(ui.LinkConfig{Href: "#allergy-note", Text: "Allergy note ↓"}))
	return Style.Scope(html.Div(html.DivConfig{Class: Style.Legend(), ExtraAttrs: html.Attrs{"role": "group", "aria-label": "Dietary labels"}},
		ui.Cluster(ui.ClusterConfig{Gap: ui.GapLG}, items...)))
}

// Render renders the jump strip, the groups and the allergy note.
func Render(cfg Config) render.HTML {
	m := cfg.Menu
	level := cfg.ItemHeadingLevel
	if level == 0 {
		level = 3
	}

	jump := make([]render.HTML, 0, len(m.Groups))
	for _, g := range m.Groups {
		jump = append(jump, ui.LinkButton(ui.LinkButtonConfig{
			Label: g.NavLabel, Href: "#" + GroupID(g.Slug) + "-items",
			Variant: ui.ButtonSecondary, Size: ui.ButtonSizeSmall,
		}))
	}

	groups := make([]render.HTML, 0, len(m.Groups))
	for gi, g := range m.Groups {
		cards := make([]render.HTML, 0, len(g.Items))
		for _, it := range g.Items {
			cards = append(cards, itemCard(it, level))
		}
		body := []render.HTML{}
		if g.Note != "" {
			body = append(body, html.Paragraph(html.TextConfig{Class: Style.Note()}, render.Text(g.Note)))
		}
		body = append(body, ui.Grid(ui.GridConfig{Min: "19rem", Gap: ui.GapMD}, cards...))
		groups = append(groups, ui.Collapsible(ui.CollapsibleConfig{
			Summary: fmt.Sprintf("%02d · %s", gi+1, g.Title),
			Open:    gi == 0,
			ID:      GroupID(g.Slug),
			// One group open at a time, the way the original menu worked:
			// the browser closes the others (exclusive <details>).
			Name: "lotus-menu",
		},
			// The jump strip targets this wrapper: following a link into
			// a closed <details> opens it, so the strip works without
			// JavaScript.
			html.Div(html.DivConfig{ID: GroupID(g.Slug) + "-items", Class: Style.Items()}, body...),
		))
	}

	return Style.Scope(html.Div(html.DivConfig{Class: Style.Menu()},
		html.Nav(html.NavConfig{Class: Style.Jump(), Label: "Menu categories"},
			ui.Cluster(ui.ClusterConfig{Gap: ui.GapSM, NoWrap: true}, jump...)),
		ui.Stack(ui.StackConfig{Gap: ui.GapSM}, groups...),
		ui.Callout(ui.CalloutConfig{Title: m.AllergyNote.Label, Variant: ui.StatusWarning, ID: "allergy-note"},
			render.Text(m.AllergyNote.Text)),
	))
}

// itemCard is one dish: the title and price in the card head, the
// description under it, and the dietary badges in the body.
func itemCard(it menu.Item, level int) render.HTML {
	var price render.HTML
	if it.Price != "" {
		price = html.Span(html.TextConfig{Class: Style.Price()}, render.Text(it.Price))
	}
	var badges render.HTML
	if len(it.Dietary) > 0 {
		bs := make([]render.HTML, 0, len(it.Dietary))
		for _, key := range it.Dietary {
			d, _ := menu.DietaryInfo(key)
			bs = append(bs, ui.Tooltip(ui.TooltipConfig{Text: d.Text},
				ui.StatusBadge(ui.StatusBadgeConfig{
					Label: d.Label, Variant: dietaryVariants[key],
					ExtraAttrs: html.Attrs{"tabindex": "0", "aria-label": d.Text},
				})))
		}
		badges = ui.Cluster(ui.ClusterConfig{Gap: ui.GapSM}, bs...)
	}
	return ui.Card(ui.CardConfig{
		Heading:      it.Title,
		HeadingLevel: level,
		Description:  it.Description,
		Action:       price,
		Variant:      ui.CardOutlined,
	}, badges)
}
