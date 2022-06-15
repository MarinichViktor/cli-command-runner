package cli

import "github.com/jroimartin/gocui"

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
