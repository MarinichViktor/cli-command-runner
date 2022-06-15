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
	input, e := getConfig()
	if e != nil {
		log.Panicln(e)
	}

	g, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	app := cli.NewAppContext(input)

	g.SetManagerFunc(cli.Layout)

	g.Update(func(g *gocui.Gui) error {
		_, e := g.SetCurrentView(cli.SERVICES_VIEW)

		if e != nil {
			return e
		}

		g.SelFgColor = gocui.ColorMagenta
		g.Highlight = true
		g.Mouse = false

		return app.SelectProject(g, app.Projects[0])
	})

	if e := app.UpdateServicesView(g); e != nil {
		log.Panicln(e)
	}

	if e := cli.SetupConsoleBindings(app, g); e != nil {
		log.Panicln(e)
	}

	if e := cli.SetupServicesBindings(app, g); e != nil {
		log.Panicln(e)
	}

	e = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for _, p := range app.Projects {
			p.Stop()
		}

		return gocui.ErrQuit
	})
	if e != nil {
		log.Panicln(e)
	}

	c := 1
	e = g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		views := []string{cli.SERVICES_VIEW, cli.CONSOLE_VIEW}
		g.SetCurrentView(views[c%len(views)])

		c++
		return nil
	})
	if e != nil {
		log.Panicln(e)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
