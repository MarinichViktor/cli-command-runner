package cli

import (
	"github.com/jroimartin/gocui"
)

//type ConsoleView struct {
//	app *Application
//}
//
//func (c *ConsoleView) Layout(g *gocui.Gui) error {
//	maxX, maxY := g.Size()
//	sView, err := g.SetView(CONSOLE_VIEW, SERVICES_W+2, 1, maxX-1, maxY-1)
//
//	if err != nil && err != gocui.ErrUnknownView {
//		return err
//	}
//
//	sView.Title = "Console"
//
//	if c.app.ActiveProject != nil {
//	}
//
//	return nil
//}

func SetupConsoleBindings(app *Application) error {
	//for _, p := range app.Projects {
	e := app.SetKeybinding(CONSOLE_VIEW, gocui.KeyCtrlA, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		view.Autoscroll = !view.Autoscroll
		return nil
	})

	if e != nil {
		return e
	}

	e = app.SetKeybinding(CONSOLE_VIEW, gocui.KeyArrowDown, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		return app.ActiveProject.View.ScrollDown(view)
		//x, y := view.Origin()
		//_, maxY := gui.Size()

		//if len(view.BufferLines())-y+1 > maxY {
		//	return view.SetOrigin(x, y+1)
		//}

		//return nil
	})

	if e != nil {
		return e
	}

	e = app.SetKeybinding(CONSOLE_VIEW, gocui.KeyArrowUp, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		return app.ActiveProject.View.ScrollUp(view)

		//x, y := view.Origin()
		//
		//if y > 0 {
		//	return view.SetOrigin(x, y-1)
		//}

		//return nil
	})

	if e != nil {
		return e
	}

	e = app.SetKeybinding(CONSOLE_VIEW, gocui.KeyArrowLeft, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, y := view.Origin()

		if x > 0 {
			return view.SetOrigin(x-1, y)
		}

		return nil
	})

	if e != nil {
		return e
	}

	e = app.SetKeybinding(CONSOLE_VIEW, gocui.KeyArrowRight, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, y := view.Origin()
		return view.SetOrigin(x+1, y)
	})

	if e != nil {
		return e
	}

	e = app.SetKeybinding(CONSOLE_VIEW, gocui.KeyPgdn, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		return app.ActiveProject.View.ScrollPageDown(view)
	})

	if e != nil {
		return e
	}

	e = app.SetKeybinding(CONSOLE_VIEW, gocui.KeyPgup, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		return app.ActiveProject.View.ScrollPageUp(view)
	})

	if e != nil {
		return e
	}

	e = app.SetKeybinding(CONSOLE_VIEW, gocui.KeyCtrlB, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, _ := view.Origin()

		view.SetOrigin(x, 0)

		return nil
	})
	if e != nil {
		return e
	}
	//}

	return nil
}
