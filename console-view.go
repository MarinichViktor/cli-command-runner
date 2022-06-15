package cli

import (
	"github.com/jroimartin/gocui"
)

func SetupConsoleBindings(app *AppContext, g *gocui.Gui) error {
	e := g.SetKeybinding(CONSOLE_VIEW, gocui.KeyCtrlA, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		view.Autoscroll = !view.Autoscroll
		return nil
	})

	if e != nil {
		return e
	}

	e = g.SetKeybinding(CONSOLE_VIEW, gocui.KeyArrowDown, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, y := view.Origin()
		_, maxY := gui.Size()

		if len(view.BufferLines())-y+1 > maxY {
			return view.SetOrigin(x, y+1)
		}

		return nil
	})

	if e != nil {
		return e
	}

	e = g.SetKeybinding(CONSOLE_VIEW, gocui.KeyArrowUp, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, y := view.Origin()

		if y > 0 {
			return view.SetOrigin(x, y-1)
		}

		return nil
	})
	e = g.SetKeybinding(CONSOLE_VIEW, gocui.KeyCtrlB, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, _ := view.Origin()

		view.SetOrigin(x, 0)

		return nil
	})
	if e != nil {
		return e
	}

	return nil
}
