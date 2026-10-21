package site

import (
	"github.com/DonaldMurillo/gofastr/core/render"
	ui "github.com/DonaldMurillo/gofastr/framework/ui"

	"github.com/AlexVCS/lotus-vegan/internal/menu"
	"github.com/AlexVCS/lotus-vegan/internal/menuview"
)

// MenuScreen is the standalone menu: the whole page is the forest band.
type MenuScreen struct{}

func (s *MenuScreen) Render() render.HTML {
	return ui.Themed(forestBand,
		ui.Container(ui.ContainerConfig{Width: ui.ContainerPage, Pad: ui.ContainerPadPage},
			ui.Stack(ui.StackConfig{Gap: ui.GapXL},
				ui.Hero(ui.HeroConfig{
					Eyebrow:  "Lotus Vegan · Full menu",
					Title:    "Find your favorite.",
					Subtitle: "Our menu is made entirely without animal products. Dietary options are marked throughout. Prices shown are from the restaurant’s current menu.",
					Footer:   menuview.Legend(),
				}),
				menuview.Render(menuview.Config{Menu: menu.Data, ItemHeadingLevel: 2}),
			),
		),
	)
}
