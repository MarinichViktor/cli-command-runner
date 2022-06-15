package cli

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

const (
	SERVICES_W    = 25
	SERVICES_VIEW = "Services"
	CONSOLE_VIEW  = "Console"
)

type ServicesView struct {
	Name     string
	Body     string
	Projects []*Project
}

func (w *ServicesView) ReDraw(v *gocui.View) {
	v.Clear()

	for _, p := range w.Projects {
		name := p.Name

		if p.IsHighlighted {
			name = fmt.Sprintf("\u001b[44m%s\033[m", name)
		}

		if p.IsRunning {
			fmt.Fprintln(v, name+" (Running)")
		} else {
			fmt.Fprintln(v, name)
		}
	}
}

func (w *ServicesView) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	v, err := g.SetView(w.Name, 1, 1, SERVICES_W+1, maxY-1)
	v.Wrap = true
	v.Title = w.Name

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		w.ReDraw(v)
	}

	return nil
}

func NewServicesView(projects []*Project, body string) *ServicesView {
	view := ServicesView{
		Name:     SERVICES_VIEW,
		Body:     body,
		Projects: projects,
	}

	return &view
}
