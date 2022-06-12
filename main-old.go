//package main
//
//import (
//	"cli/command"
//	"fmt"
//	"github.com/jroimartin/gocui"
//	"log"
//)
//
//func main() {
//	g, err := gocui.NewGui(gocui.Output256)
//	if err != nil {
//		log.Panicln(err)
//	}
//	defer g.Close()
//
//	g.Cursor = true
//	g.Mouse = true
//	g.SetManager()
//	g.SetManagerFunc(layout)
//
//	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
//		log.Panicln(err)
//	}
//
//	g.Update(func(gui *gocui.Gui) error {
//		v, _ := g.SetCurrentView("out")
//		v.Title = "Console"
//
//		return nil
//	})
//	activeItem := 0
//	items := []string{
//		"tas",
//	}
//	e := g.SetKeybinding("hello", gocui.KeyArrowDown, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
//		//view.Clear()
//		//view.Title = "321"
//		if activeItem == len(items)-1 {
//			activeItem = 0
//		} else {
//			activeItem++
//		}
//		gui.Update(func(gui *gocui.Gui) error {
//			view.Clear()
//
//			for dx, i := range items {
//				if activeItem == dx {
//					fmt.Fprintf(view, "\u001b[44m%s\033[m\n", i)
//				} else {
//					fmt.Fprintln(view, i)
//				}
//			}
//			return nil
//		})
//
//		return nil
//	})
//	if e != nil {
//		log.Panicln(err)
//	}
//
//	cmdRunner, e := command.NewCommandRunner("docker-compose up", "/home/vmaryn/projects/go/sandbox")
//
//	if e != nil {
//		panic(e)
//	}
//	e = cmdRunner.Start()
//
//	if e != nil {
//		panic(e)
//	}
//	go func() {
//		for {
//			select {
//			case v, ok := <-cmdRunner.OutStream:
//				if !ok {
//					cmdRunner.Stop()
//					return
//				}
//
//				g.Update(func(gui *gocui.Gui) error {
//					g.SetCurrentView("out")
//					view, _ := g.View("out")
//					fmt.Fprint(view, v)
//					return nil
//				})
//			}
//		}
//
//	}()
//
//	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
//		log.Panicln(err)
//	}
//}
//
//func layout(g *gocui.Gui) error {
//	maxX, maxY := g.Size()
//	var activeItem = 0
//
//	if v, err := g.SetView("services", 0, 0, 20, maxY/2); err != nil {
//		if err != gocui.ErrUnknownView {
//			return err
//		}
//
//		items := []string{
//			"tas", "admin", "mobile-api",
//		}
//		for dx, i := range items {
//			if activeItem == dx {
//				fmt.Fprintf(v, "\u001b[44m%s\033[m\n", i)
//			} else {
//				fmt.Fprintln(v, i)
//			}
//		}
//	}
//
//	if v, err := g.SetView("out", 21, 0, maxX-1, maxY-1); err != nil {
//		if err != gocui.ErrUnknownView {
//			return err
//		}
//		//
//		//items := []string{
//		//	"tas", "admin", "mobile-api",
//		//}
//		//for dx, i := range items {
//		//	if activeItem == dx {
//		//		fmt.Fprintf(v, "\u001b[44m%s\033[m\n", i)
//		//	} else {
//		//		fmt.Fprintln(v, i)
//		//	}
//		//}
//		v.Autoscroll = true
//		v.Wrap = true
//	}
//	return nil
//}
//
//func quit(g *gocui.Gui, v *gocui.View) error {
//	return gocui.ErrQuit
//}
