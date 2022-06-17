package cli

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	sView, err := g.SetView(CONSOLE_VIEW, SERVICES_W+2, 1, maxX-1, maxY-1)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		sView.Autoscroll = true
		sView.Title = "Console"
	}

	cView, err := g.SetView(SERVICES_VIEW, 1, 1, SERVICES_W+1, maxY-1)
	cView.Wrap = true

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		cView.Title = "Services"
	}

	return nil
}

func LayoutFactory(app *Application) func(g *gocui.Gui) error {
	i := 0
	return func(g *gocui.Gui) error {
		maxX, maxY := app.Size()
		for _, p := range app.Projects {
			view, err := app.SetView(p.Name, SERVICES_W+2, 1, maxX-1, maxY-1)

			if err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}

				view.Autoscroll = true
				//view.Wrap = true
				view.Title = fmt.Sprintf("Console - %s - %d", p.Name, i)
			}
		}

		cView, err := g.SetView(SERVICES_VIEW, 1, 1, SERVICES_W+1, maxY-1)
		cView.Wrap = true
		i++

		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

		}
		cView.Title = fmt.Sprintf("Services - %d", i)

		return nil
	}
}
