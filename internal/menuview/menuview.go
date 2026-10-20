// Package menuview renders the Lotus menu (internal/menu) with GoFastr
// components: a Collapsible per menu group (any number open), a Grid of
// outlined Cards for the items (title, price, description, a list of
// kinds, dietary badges), StatusBadge + Tooltip for the dietary codes, a jump strip of
// LinkButtons, and a Callout for the allergy note.
//
// The few things the kit does not draw (the scrolling jump strip, the
// price colour, the divided list of kinds) live in menuview.style.css, an owned
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
	return Style.Scope(html.Div(html.DivConfig{Class: Style.RootWith(RootVariants{Legend: true}), ExtraAttrs: html.Attrs{"role": "group", "aria-label": "Dietary labels"}},
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
			// No Name: groups open independently, so several can be
			// open at once.
		},
			// The jump strip targets this wrapper: following a link into
			// a closed <details> opens it, so the strip works without
			// JavaScript.
			html.Div(html.DivConfig{ID: GroupID(g.Slug) + "-items", Class: Style.Items()}, body...),
		))
	}

	return Style.Scope(html.Div(html.DivConfig{Class: Style.RootWith(RootVariants{Menu: true})},
		html.Nav(html.NavConfig{Class: Style.Jump(), Label: "Menu categories"},
			ui.Cluster(ui.ClusterConfig{Gap: ui.GapSM, NoWrap: true}, jump...)),
		ui.Stack(ui.StackConfig{Gap: ui.GapSM}, groups...),
		// The anchor and its spacing live on a wrapper we own: an owned
		// style may only place a kit root, and scroll-margin restyles it.
		html.Div(html.DivConfig{ID: "allergy-note", Class: Style.Allergy()},
			ui.Callout(ui.CalloutConfig{Title: m.AllergyNote.Label, Variant: ui.StatusWarning},
				render.Text(m.AllergyNote.Text))),
	))
}

// itemCard is one dish: the title and price in the card head, the
// description under it, and in the body the kinds (one divided row
// each: name, then ingredients) and the dietary badges.
func itemCard(it menu.Item, level int) render.HTML {
	var price render.HTML
	if it.Price != "" {
		price = html.Span(html.TextConfig{Class: Style.Price()}, render.Text(it.Price))
	}
	body := []render.HTML{}
	if len(it.Kinds) > 0 {
		rows := make([]render.HTML, 0, len(it.Kinds))
		for _, k := range it.Kinds {
			rows = append(rows, html.ListItem(html.ListItemConfig{Class: Style.KindRow()},
				html.Span(html.TextConfig{Class: Style.KindName()}, render.Text(k.Name)),
				html.Span(html.TextConfig{Class: Style.KindIngredients()}, render.Text(k.Ingredients)),
			))
		}
		body = append(body, html.UnorderedList(html.ListConfig{Class: Style.Kinds()}, rows...))
	}
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
		body = append(body, ui.Cluster(ui.ClusterConfig{Gap: ui.GapSM}, bs...))
	}
	var content render.HTML
	switch len(body) {
	case 1:
		content = body[0]
	case 2:
		content = ui.Stack(ui.StackConfig{Gap: ui.GapLG}, body...)
	}
	return ui.Card(ui.CardConfig{
		Heading:      it.Title,
		HeadingLevel: level,
		Description:  it.Description,
		Action:       price,
		Variant:      ui.CardOutlined,
	}, content)
}
