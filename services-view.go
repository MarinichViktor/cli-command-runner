package cli

import (
	"github.com/jroimartin/gocui"
	"log"
)

const (
	SERVICES_W    = 25
	SERVICES_VIEW = "Services"
	CONSOLE_VIEW  = "Console"
)

func SetupServicesBindings(app *AppContext, g *gocui.Gui) error {
	e := g.SetKeybinding(SERVICES_VIEW, gocui.KeyArrowDown, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for i, p := range app.Projects {
			if p.IsHighlighted && i < len(app.Projects)-1 {
				p.IsHighlighted = false
				app.Projects[i+1].IsHighlighted = true
				if e := app.SelectProject(g, app.Projects[i+1]); e != nil {
					log.Panicln(e)
				}

				if e := app.UpdateServicesView(g); e != nil {
					return e
				}

				break
			}
		}

		return nil
	})

	if e != nil {
		return e
	}

	e = g.SetKeybinding(SERVICES_VIEW, gocui.KeyArrowUp, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for i, p := range app.Projects {
			if p.IsHighlighted && i != 0 {
				g.Update(func(gui *gocui.Gui) error {
					p.IsHighlighted = false
					app.Projects[i-1].IsHighlighted = true
					if e := app.SelectProject(g, app.Projects[i-1]); e != nil {
						log.Panicln(e)
					}

					app.UpdateServicesView(g)
					return nil
				})

				break
			}
		}

		return nil
	})

	if e != nil {
		return e
	}

	e = g.SetKeybinding(SERVICES_VIEW, gocui.KeyCtrlR, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for _, p := range app.Projects {
			if p.IsHighlighted {
				if p.IsRunning {
					// todo: to be fixed
					if e := p.Stop(); e != nil {
						return e
					}
				} else {
					if e := p.Start(); e != nil {
						return e
					}

					if e := app.UpdateServicesView(g); e != nil {
						return e
					}
				}

				break
			}
		}
		return nil
	})

	return e
}
