package main

import (
	"cli"
	"github.com/jroimartin/gocui"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

func getConfig() ([]*cli.ProjectArgs, error) {
	l, _ := strconv.Atoi(os.Args[2])
	args := make([]*cli.ProjectArgs, l)

	data, e := os.ReadFile(os.Args[1])
	if e != nil {
		return nil, e
	}

	if e := yaml.Unmarshal(data, &args); e != nil {
		return nil, e
	}

	return args, nil
}

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

	c := 1
	e = app.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		views := []string{cli.SERVICES_VIEW}
		for _, p := range app.Projects {
			if p.IsHighlighted {
				views = append(views, p.Name)
				break
			}
		}
		app.SetCurrentView(views[c%len(views)])

		c++
		return nil
	})

	return e
}
