package cli

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const APP_KEY = "app"

type AppContext struct {
	g            *gocui.Gui
	Active       *Project
	Projects     []*Project
	Services     *ServicesView
	Console      *ConsoleView
	ConBuff      *ConsoleBuff
	done         chan struct{}
	unsubscriber func()
}

func NewAppContext(projArgs []*ProjectArgs, g *gocui.Gui) *AppContext {
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
		g:        g,
	}
}

func (app *AppContext) ServicesView() (*gocui.View, error) {
	v, e := app.g.View(SERVICES_VIEW)
	if e != nil && e != gocui.ErrUnknownView {
		return nil, e
	}

	return v, nil
}

func (app *AppContext) ConsoleView() (*gocui.View, error) {
	v, e := app.g.View(CONSOLE_VIEW)
	if e != nil && e != gocui.ErrUnknownView {
		return nil, e
	}

	return v, nil
}

func (app *AppContext) UpdateServicesView() error {
	app.g.Update(func(gui *gocui.Gui) error {
		view, e := app.ServicesView()

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

func (app *AppContext) UpdateConsoleView(data string) error {
	app.g.Update(func(gui *gocui.Gui) error {
		view, e := app.ConsoleView()
		if e != nil {
			return e
		}
		view.Clear()
		_, e = fmt.Fprintln(view, data)

		return e
	})

	return nil
}

func (app *AppContext) SelectProject(p *Project) error {
	if app.unsubscriber != nil {
		app.unsubscriber()
	}

	app.unsubscriber = p.Subscribe(func(data string) {
		if e := app.UpdateConsoleView(p.Data); e != nil {
			return
		}
	}, func() {
		app.unsubscriber()
		if e := app.UpdateServicesView(); e != nil {
			return
		}
	})

	return app.UpdateConsoleView(p.Data)
}

type ConsoleBuff struct {
	Project *Project
	Data    []string
}
