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

	app := cli.NewAppContext(input, g)

	g.SetManagerFunc(layout)

	g.Update(func(gui *gocui.Gui) error {
		_, e := g.SetCurrentView(cli.SERVICES_VIEW)
		g.SelBgColor = gocui.ColorWhite
		g.SelFgColor = gocui.ColorBlue
		g.Highlight = true
		return e
	})

	if e := app.UpdateServicesView(); e != nil {
		panic(e)
	}

	if e := cli.SetupConsoleBindings(g); e != nil {
		panic(e)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		//cmdRunner.Stop()
		return gocui.ErrQuit
	}); err != nil {
		log.Panicln(err)
	}

	c := 1
	g.SetKeybinding("", gocui.KeyCtrl2, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		views := []string{cli.SERVICES_VIEW, cli.CONSOLE_VIEW}
		g.SetCurrentView(views[c%len(views)])

		c++
		return nil
	})

	g.SetKeybinding(cli.SERVICES_VIEW, gocui.KeyArrowDown, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for i, p := range app.Projects {
			if p.IsHighlighted && i < len(app.Projects)-1 {
				p.IsHighlighted = false
				app.Projects[i+1].IsHighlighted = true

				if e := app.UpdateServicesView(); e != nil {
					return e
				}

				break
			}
		}

		return nil
	})

	g.SetKeybinding(cli.SERVICES_VIEW, gocui.KeyArrowUp, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for i, p := range app.Projects {
			if p.IsHighlighted && i != 0 {
				g.Update(func(gui *gocui.Gui) error {
					p.IsHighlighted = false
					app.Projects[i-1].IsHighlighted = true
					app.UpdateServicesView()
					return nil
				})

				break
			}
		}

		return nil
	})

	g.SetKeybinding(cli.SERVICES_VIEW, gocui.KeyCtrlR, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for _, p := range app.Projects {
			if p.IsHighlighted {
				if p.IsRunning {
					p.CmdInst.Stop()
				} else {
					if e := p.Start(); e != nil {
						return e
					}

					if e := app.UpdateServicesView(); e != nil {
						return e
					}
				}

				break
			}
		}
		return nil
	})

	//g.Update(func(gui *gocui.Gui) error {
	//	_, e := g.SetCurrentView(cli.SERVICES_VIEW)
	//	if e != nil {
	//		log.Panicln(e)
	//	}
	//
	//	v, _ := g.View(cli.CONSOLE_VIEW)
	//	v.Editor = gocui.DefaultEditor
	//	v.Editable = true
	//	g.SelBgColor = gocui.ColorWhite
	//	g.SelFgColor = gocui.ColorBlue
	//	g.Highlight = true
	//	return nil
	//})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	sView, err := g.SetView(cli.CONSOLE_VIEW, cli.SERVICES_W+2, 1, maxX-1, maxY-1)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		sView.Title = "Console"
	}

	cView, err := g.SetView(cli.SERVICES_VIEW, 1, 1, cli.SERVICES_W+1, maxY-1)
	cView.Wrap = true

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		cView.Title = "Services"
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
