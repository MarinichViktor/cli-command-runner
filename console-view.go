package cli

import (
	"github.com/jroimartin/gocui"
)

func SetupConsoleBindings(app *Application) error {
	for _, p := range app.Projects {
		e := app.SetKeybinding(p.Name, gocui.KeyCtrlA, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
			view.Autoscroll = !view.Autoscroll
			return nil
		})

		if e != nil {
			return e
		}

		e = app.SetKeybinding(p.Name, gocui.KeyArrowDown, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
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

		e = app.SetKeybinding(p.Name, gocui.KeyArrowUp, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
			x, y := view.Origin()

			if y > 0 {
				return view.SetOrigin(x, y-1)
			}

			return nil
		})
		e = app.SetKeybinding(p.Name, gocui.KeyCtrlB, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
			x, _ := view.Origin()

			view.SetOrigin(x, 0)

			return nil
		})
		if e != nil {
			return e
		}
	}

	return nil
}
