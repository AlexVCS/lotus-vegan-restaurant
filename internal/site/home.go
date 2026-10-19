package site

import (
	"github.com/DonaldMurillo/gofastr/core-ui/html"
	"github.com/DonaldMurillo/gofastr/core/render"
	ui "github.com/DonaldMurillo/gofastr/framework/ui"

	"github.com/AlexVCS/lotus-vegan/internal/menu"
	"github.com/AlexVCS/lotus-vegan/internal/menuview"
)

// HomeScreen is the one-page site: hero with the pickup phones, the
// quick facts, the food story and gallery, the pizza band, the full menu
// and the visit details.
type HomeScreen struct{}

func (s *HomeScreen) Render() render.HTML {
	return ui.Stack(ui.StackConfig{Gap: ui.GapNone},
		ui.Container(ui.ContainerConfig{Width: ui.ContainerPage, Pad: ui.ContainerPadPage}, hero()),
		ui.Themed(limeBand, ui.Container(ui.ContainerConfig{Width: ui.ContainerPage, Pad: ui.ContainerPadPage}, quickFacts())),
		ui.Container(ui.ContainerConfig{Width: ui.ContainerPage, Pad: ui.ContainerPadPage}, story()),
		ui.Themed(goldBand, ui.Container(ui.ContainerConfig{Width: ui.ContainerPage, Pad: ui.ContainerPadPage}, pizza())),
		ui.Themed(forestBand, ui.Container(ui.ContainerConfig{Width: ui.ContainerPage, Pad: ui.ContainerPadPage},
			ui.Section(ui.SectionConfig{
				ID:          "menu",
				Eyebrow:     "Full menu",
				Heading:     "Find your favorite.",
				Description: "Our menu is made entirely without animal products. Dietary options are marked throughout.",
			},
				ui.Stack(ui.StackConfig{Gap: ui.GapXL},
					menuview.Legend(),
					menuview.Render(menuview.Config{Menu: menu.Data, ItemHeadingLevel: 3}),
				),
			))),
		ui.Themed(limeBand, ui.Container(ui.ContainerConfig{Width: ui.ContainerPage, Pad: ui.ContainerPadPage}, visit())),
	)
}

func hero() render.HTML {
	actions := make([]render.HTML, 0, len(phones)+1)
	for i, p := range phones {
		variant := ui.ButtonPrimary
		if i > 0 {
			variant = ui.ButtonSecondary
		}
		actions = append(actions, ui.LinkButton(ui.LinkButtonConfig{
			Label: "Call " + p.Display, Href: p.Tel, Icon: "phone", Variant: variant, Size: ui.ButtonSizeLarge,
			ExtraAttrs: html.Attrs{"aria-label": "Call " + p.Display + " to order pickup"},
		}))
	}
	actions = append(actions, ui.LinkButton(ui.LinkButtonConfig{
		Label: "Explore the menu ↘", Href: "#menu", Variant: ui.ButtonGhost, Size: ui.ButtonSizeLarge,
	}))

	media := html.Figure(html.FigureConfig{Class: HomeStyle.Media()},
		ui.OptimizedImage(ui.OptimizedImageConfig{
			Src:          "/assets/avocado-salad.jpg",
			Alt:          "Lotus Vegan avocado salad with spinach, tomato, cucumber, onion and plant-based protein",
			Width:        1067,
			Height:       800,
			Eager:        true,
			HighPriority: true,
			Rounded:      true,
		}),
		html.Span(html.TextConfig{Class: HomeStyle.Stamp(), ExtraAttrs: html.Attrs{"aria-hidden": "true"}},
			html.Strong(html.TextConfig{}, render.Text("100%")), render.Text(" vegan")),
		html.FigCaption(html.FigCaptionConfig{Class: HomeStyle.Caption()},
			html.Span(html.TextConfig{}, render.Text("Avocado Salad")),
			html.Span(html.TextConfig{Class: HomeStyle.CaptionNo()}, render.Text("No. 19"))),
	)

	return ui.Hero(ui.HeroConfig{
		Eyebrow:  "Nhà Hàng Chay Hoa Sen · Pinellas Park",
		Title:    "Food with peace at heart.",
		Subtitle: "A joyful, fully vegan kitchen serving Vietnamese tradition, raw favorites, noodle bowls, fresh juices, and abundant comfort food.",
		Actions:  actions,
		Footer: html.Paragraph(html.TextConfig{Class: HomeStyle.Pickup()},
			html.Strong(html.TextConfig{}, render.Text("Order pickup now")),
			render.Text(" — call either number.")),
		Media: media,
	})
}

func quickFacts() render.HTML {
	return ui.Grid(ui.GridConfig{Min: "15rem", Gap: ui.GapLG, ExtraAttrs: html.Attrs{"role": "group", "aria-label": "Restaurant details"}},
		ui.FactBox(ui.FactBoxConfig{Label: "Open Mon–Sat", Value: "11am–9pm · last dine-in 8:15pm"}),
		ui.FactBox(ui.FactBoxConfig{Label: "Second Sunday buffet", Value: "12pm–7pm · all other Sundays closed"}),
		ui.FactBox(ui.FactBoxConfig{Label: "Pinellas Park", Value: "6575 Park Blvd N, Florida 33781"}),
	)
}

func story() render.HTML {
	cards := make([]render.HTML, 0, len(dishes))
	for _, d := range dishes {
		cards = append(cards, ui.Card(ui.CardConfig{
			Href:    d.Href,
			Variant: ui.CardOutlined,
			Header: ui.OptimizedImage(ui.OptimizedImageConfig{
				Src: d.Src, Alt: d.Alt, Width: d.Width, Height: d.Height,
				Aspect: ui.ImageAspectSquare, Fit: ui.ImageFitCover,
			}),
			ExtraAttrs: html.Attrs{"target": "_blank", "rel": "noreferrer", "aria-label": d.Label},
		},
			html.Strong(html.TextConfig{Class: HomeStyle.DishName()}, render.Text(d.Name)),
			ui.Muted(render.Text(d.Note)),
		))
	}

	return ui.Section(ui.SectionConfig{
		ID:          "story",
		Eyebrow:     "Our food",
		Heading:     "Be vegan, make peace.",
		Description: "Lotus Vegan brings a generous plant-based table to the heart of Pinellas Park—rooted in Vietnamese flavors and made for everyone.",
	},
		ui.Stack(ui.StackConfig{Gap: ui.Gap2XL},
			ui.Stack(ui.StackConfig{Gap: ui.GapLG},
				html.Paragraph(html.TextConfig{Class: HomeStyle.Lead()}, render.Text(
					"From fragrant phở and crisp noodles to raw meals, hand-tossed pizza and boba, the menu offers a variety of delightful vegan dishes. Gluten-free and soy-free options are marked throughout the menu.")),
				ui.Cluster(ui.ClusterConfig{Gap: ui.GapMD},
					ui.LinkButton(ui.LinkButtonConfig{Label: "See menu →", Href: "#menu"}),
					ui.LinkButton(ui.LinkButtonConfig{
						Label: "Follow Lotus on Instagram", Href: instagramURL, Icon: "instagram", External: true, Variant: ui.ButtonSecondary,
						ExtraAttrs: html.Attrs{"aria-label": "Follow Lotus Vegan on Instagram (opens in a new tab)"},
					}),
					ui.LinkButton(ui.LinkButtonConfig{
						Label: "Follow Lotus on Facebook", Href: facebookURL, Icon: "facebook", External: true, Variant: ui.ButtonSecondary,
						ExtraAttrs: html.Attrs{"aria-label": "Follow Lotus Vegan on Facebook (opens in a new tab)"},
					}),
				),
			),
			ui.Grid(ui.GridConfig{Min: "15rem", Gap: ui.GapLG}, cards...),
		),
	)
}

func pizza() render.HTML {
	return ui.Section(ui.SectionConfig{
		Class:       HomeStyle.PizzaSection(),
		Eyebrow:     "Now from the oven",
		Heading:     "Vegan pizza, the Lotus way.",
		Description: "Housemade dough, marinara, vegan cheese, bell peppers, vegan pepperoni, mushrooms, and onion. Served with a side salad.",
	},
		html.Div(html.DivConfig{Class: HomeStyle.PizzaPrices(), ExtraAttrs: html.Attrs{"role": "group", "aria-label": "Pizza prices"}},
			html.Div(html.DivConfig{Class: HomeStyle.PizzaSizes()},
				pizzaSize("9 inch", "$18"),
				pizzaSize("16 inch", "$26"),
			),
			html.Paragraph(html.TextConfig{Class: HomeStyle.PizzaExtra()},
				render.Text("Extra toppings "), html.Strong(html.TextConfig{}, render.Text("+$2"))),
		),
	)
}

// pizzaSize is one size tile in the pizza band: the size on top, the
// price under it.
func pizzaSize(size, price string) render.HTML {
	return html.Div(html.DivConfig{Class: HomeStyle.PizzaSize()},
		html.Span(html.TextConfig{Class: HomeStyle.PizzaSizeLabel()}, render.Text(size)),
		html.Span(html.TextConfig{Class: HomeStyle.PizzaSizePrice()}, render.Text(price)),
	)
}

// visitFact is one column of the visit facts: a heading, the value, and
// an optional footer under a divider. The columns share their rows (see
// .visit-fact), so the dividers sit on one line even when one footer
// runs taller than another. ui.Card can't do that: it pins its footer
// to the card's bottom, which pushes a shorter footer down a shared row.
func visitFact(heading string, value, footer render.HTML) render.HTML {
	parts := []render.HTML{
		html.Heading(html.HeadingConfig{Level: 3, Class: HomeStyle.VisitHeading()}, render.Text(heading)),
		html.Div(html.DivConfig{Class: HomeStyle.VisitBody()}, value),
	}
	if footer != "" {
		parts = append(parts, html.Div(html.DivConfig{Class: HomeStyle.VisitFoot()}, footer))
	}
	return html.Div(html.DivConfig{Class: HomeStyle.VisitFact()}, parts...)
}

func visit() render.HTML {
	callLinks := make([]render.HTML, 0, len(phones))
	for _, p := range phones {
		callLinks = append(callLinks, ui.LinkButton(ui.LinkButtonConfig{
			Label: p.Display, Href: p.Tel, Icon: "phone", Variant: ui.ButtonSecondary,
			ExtraAttrs: html.Attrs{"aria-label": "Call " + p.Display},
		}))
	}
	return ui.Section(ui.SectionConfig{
		ID:      "visit",
		Eyebrow: "Come sit with us",
		Heading: "Your table is waiting.",
	},
		html.Div(html.DivConfig{Class: HomeStyle.VisitFacts()},
			visitFact("Address",
				html.Paragraph(html.TextConfig{Class: HomeStyle.VisitValue()}, render.Text("6575 Park Blvd N"), html.BR(), render.Text("Pinellas Park, FL 33781")),
				ui.LinkButton(ui.LinkButtonConfig{
					Label: "Get directions", Href: directionsURL, Icon: "map-pin", External: true, Variant: ui.ButtonPrimary, Size: ui.ButtonSizeSmall,
				}),
			),
			visitFact("Hours",
				html.Paragraph(html.TextConfig{Class: HomeStyle.VisitValue()}, render.Text("Monday–Saturday"), html.BR(), render.Text("11:00am–9:00pm")),
				ui.Muted(render.Text("Last dine-in at 8:15pm · Second Sunday buffet 12pm–7pm")),
			),
			visitFact("Call",
				ui.Stack(ui.StackConfig{Gap: ui.GapSM, Align: ui.AlignStart}, callLinks...),
				"",
			),
		),
	)
}
