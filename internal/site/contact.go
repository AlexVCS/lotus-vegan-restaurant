package site

import (
	"github.com/DonaldMurillo/gofastr/core-ui/html"
	"github.com/DonaldMurillo/gofastr/core/render"
	ui "github.com/DonaldMurillo/gofastr/framework/ui"
)

// ContactScreen is the contact page: the note on the left, a forest card
// on the right with the ways to reach the restaurant. There is no online
// form: nothing would receive the message, so the card sends visitors to
// the phones instead.
type ContactScreen struct{}

func (s *ContactScreen) Render() render.HTML {
	intro := ui.Stack(ui.StackConfig{Gap: ui.GapLG},
		ui.Hero(ui.HeroConfig{
			Eyebrow:  "Contact us",
			Title:    "Talk to the Lotus team.",
			Subtitle: "Questions about the menu, dietary needs, or your next visit? Give us a call and we’ll be happy to help.",
		}),
	)

	// The card is drawn in the forest band's palette; the frame clips the
	// band's square backdrop to the card's corners.
	card := html.Div(html.DivConfig{Class: ContactStyle.Frame()},
		ui.Themed(forestBand, ui.Card(ui.CardConfig{Heading: "Say hello", HeadingLevel: 2, ID: "contact"},
			ui.Stack(ui.StackConfig{Gap: ui.GapLG},
				html.Paragraph(html.TextConfig{}, render.Text(
					"We can’t take messages online just yet. Call either number and the Lotus team will help with menu questions, allergies, or pickup orders.")),
				callLinks(),
				ui.Muted(ui.Link(ui.LinkConfig{Href: "/#visit", Text: "See our hours and address"})),
			))))

	return ui.Container(ui.ContainerConfig{Width: ui.ContainerPage, Pad: ui.ContainerPadPage},
		ui.Grid(ui.GridConfig{Min: "20rem", Gap: ui.Gap2XL}, intro, card))
}
