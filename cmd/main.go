package main

import (
	"cli"
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	app, e := cli.NewApp()

	if e != nil {
		log.Panicln(e)
	}

	defer app.Close()

	app.SetManagerFunc(cli.LayoutFactory(app))
	app.Update(func(g *gocui.Gui) error {
		_, e := g.SetCurrentView(cli.SERVICES_VIEW)

		if e != nil {
			return e
		}

		g.SelFgColor = gocui.ColorMagenta
		g.Highlight = true
		g.Mouse = false

		return app.SelectProject(app.Projects[0])
	})

	if e := app.UpdateServicesView(); e != nil {
		log.Panicln(e)
	}

	if e := SetupBindings(app); e != nil {
		log.Panicln(e)
	}

	if err := app.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func SetupBindings(app *cli.Application) error {
	if e := cli.SetupConsoleBindings(app); e != nil {
		log.Panicln(e)
	}

	if e := cli.SetupServicesBindings(app); e != nil {
		log.Panicln(e)
	}

	e := app.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for _, p := range app.Projects {
			p.Stop()
		}

		return gocui.ErrQuit
	})
	if e != nil {
		log.Panicln(e)
	}

	// todo: replace this ugly solution
	windowIdx := 1
	e = app.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		views := []string{cli.SERVICES_VIEW, cli.CONSOLE_VIEW}
		app.SetCurrentView(views[windowIdx%len(views)])

		windowIdx++
		return nil
	})

	return e
}
