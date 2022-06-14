package cli

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type ConsoleView struct {
	Name string
	Body string
}

func (w *ConsoleView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(w.Name, SERVICES_W+2, 1, maxX-1, maxY-1)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Title = w.Name
		fmt.Fprintln(v, w.Body)
		lines := v.BufferLines()
		if len(lines) > 0 {
			v.SetCursor(len(v.BufferLines()), len(lines[len(lines)-1]))
		}
	}

	return nil
}

func NewConsoleView(body string) *ConsoleView {
	consoleView := ConsoleView{
		Name: CONSOLE_VIEW,
		Body: "",
	}

	return &consoleView
}

func SetupConsoleBindings(g *gocui.Gui) error {
	e := g.SetKeybinding(CONSOLE_VIEW, gocui.KeyCtrlA, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		view.Autoscroll = !view.Autoscroll
		return nil
	})

	if e != nil {
		return e
	}

	e = g.SetKeybinding(CONSOLE_VIEW, gocui.KeyArrowDown, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, y := view.Origin()
		_, maxY := gui.Size()

		if len(view.BufferLines())-y+1 > maxY {
			return view.SetOrigin(x, y+1)
		}

		return nil
	})

	if e != nil {
		return e
	}

	e = g.SetKeybinding(CONSOLE_VIEW, gocui.KeyArrowUp, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, y := view.Origin()

		if y > 0 {
			return view.SetOrigin(x, y-1)
		}

		return nil
	})

	if e != nil {
		return e
	}

	return nil
}
