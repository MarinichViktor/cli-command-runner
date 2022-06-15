package cli

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type Application struct {
	*gocui.Gui
	Projects []*Project
}

func (app *Application) SelectProject(p *Project) error {
	maxX, maxY := app.Size()
	view, e := app.SetView(p.Name, SERVICES_W+2, 1, maxX-1, maxY-1)
	view.Wrap = true

	if e != nil && e != gocui.ErrUnknownView {
		return e
	}

	if _, e := app.SetViewOnTop(p.Name); e != nil {
		return e
	}

	p.Subscribe(func(data string) {
		if !p.IsHighlighted {
			return
		}

		app.Update(func(gui *gocui.Gui) error {
			view.Clear()
			if _, e := fmt.Fprint(view, data); e != nil {
				return e
			}

			return nil
		})
	}, func() {
		app.UpdateServicesView()
		if p.IsHighlighted {
			app.Update(func(gui *gocui.Gui) error {
				p.Data = append(p.Data, "Command exited....")
				if _, e := fmt.Fprintln(view, "Command exited...."); e != nil {
					return e
				}

				return nil
			})
		}
	})

	app.Update(func(gui *gocui.Gui) error {
		view.Clear()
		if _, e := fmt.Fprint(view, p.StrData()); e != nil {
			return e
		}

		return nil
	})

	return nil
}
func (app *Application) UpdateServicesView() error {
	app.Update(func(gui *gocui.Gui) error {
		view, e := gui.View(SERVICES_VIEW)
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

func NewApp() (*Application, error) {
	g, e := gocui.NewGui(gocui.Output256)

	if e != nil {
		return nil, e
	}
	projects, e := buildProjects()

	if e != nil {
		return nil, e
	}

	return &Application{
		Gui:      g,
		Projects: projects,
	}, nil
}

func buildProjects() ([]*Project, error) {
	args, e := ParseProjectArgs()
	projects := make([]*Project, 0)

	if e != nil {
		return nil, e
	}

	for i, a := range args {
		project := NewProject(a)

		if i == 0 {
			project.IsHighlighted = true
		}

		projects = append(projects, project)
	}

	return projects, nil
}
