package cli

//
//import (
//	"fmt"
//	"github.com/jroimartin/gocui"
//)
//
//type ViewBox struct {
//	Name string
//	Body string
//	X0   func(int) int
//	Y0   func(int) int
//	X1   func(int) int
//	Y1   func(int) int
//}
//
//func (w *ViewBox) Layout(g *gocui.Gui) error {
//	maxX, maxY := g.Size()
//	v, err := g.SetView(w.Name, w.X0(maxX), w.Y0(maxY), w.X1(maxX), w.Y1(maxY))
//	v.Wrap = true
//	v.Title = w.Name
//
//	if err != nil {
//		if err != gocui.ErrUnknownView {
//			return err
//		}
//
//		fmt.Fprintln(v, w.Body)
//		lines := v.BufferLines()
//		if len(lines) > 0 {
//			v.SetCursor(len(v.BufferLines()), len(lines[len(lines)-1]))
//		}
//	}
//
//	return nil
//}
