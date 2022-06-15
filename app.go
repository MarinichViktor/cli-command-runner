package cli

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type AppContext struct {
	Projects     []*Project
	unsubscriber func()
}

func NewAppContext(projArgs []*ProjectArgs) *AppContext {
	projects := make([]*Project, 0)

	for i, a := range projArgs {
		project := NewProject(a)

		if i == 0 {
			project.IsHighlighted = true
		}

		projects = append(projects, project)
	}

	return &AppContext{
		Projects: projects,
	}
}

func (app *AppContext) ServicesView(g *gocui.Gui) (*gocui.View, error) {
	v, e := g.View(SERVICES_VIEW)
	if e != nil && e != gocui.ErrUnknownView {
		return nil, e
	}

	return v, nil
}

func (app *AppContext) ConsoleView(g *gocui.Gui) (*gocui.View, error) {
	v, e := g.View(CONSOLE_VIEW)
	if e != nil && e != gocui.ErrUnknownView {
		return nil, e
	}

	return v, nil
}

func (app *AppContext) UpdateServicesView(g *gocui.Gui) error {
	g.Update(func(gui *gocui.Gui) error {
		view, e := app.ServicesView(g)

		if e != nil {
			return e
		}
		view.Clear()
		for _, p := range app.Projects {
			name := p.Name

			if p.IsHighlighted {
				name = fmt.Sprintf("\u001b[44m%s\033[m", name)
			}

			if p.IsRunning {
				_, e = fmt.Fprintln(view, name+" (Running)")
				if e != nil {
					return e
				}
			} else {
				_, e = fmt.Fprintln(view, name)
				if e != nil {
					return e
				}
			}
		}

		return nil
	})

	return nil
}

func (app *AppContext) UpdateConsoleView(g *gocui.Gui, data string) error {
	g.Update(func(gui *gocui.Gui) error {
		view, e := app.ConsoleView(g)
		if e != nil {
			return e
		}
		view.Clear()
		_, e = fmt.Fprintln(view, data)

		return e
	})

	return nil
}

func (app *AppContext) SelectProject(g *gocui.Gui, p *Project) error {
	if app.unsubscriber != nil {
		app.unsubscriber()
	}
	view, e := app.ConsoleView(g)
	if e != nil {
		return e
	}

	x, _ := view.Origin()
	e = view.SetOrigin(x, 0)
	if e != nil {
		return e
	}

	app.unsubscriber = p.Subscribe(func(data string) {
		if e := app.UpdateConsoleView(g, p.StrData()); e != nil {
			return
		}
	}, func() {
		if e := app.UpdateServicesView(g); e != nil {
			return
		}
		app.unsubscriber()
	})

	return app.UpdateConsoleView(g, p.StrData())
}
