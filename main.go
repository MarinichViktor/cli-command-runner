package main

import (
	"cli/command"
	"cli/views"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

//func main() {
//	cmd := exec.Command("docker", "exec", "-it", "af57f426e5f5", "bash")
//	//cmd.Stdin = os.Stdin
//	cmd.Stdout = os.Stdout
//	//cmd.Stderr = os.Stderr
//	cmd.Start()
//	cmd.Wait()
//}

func main() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true
	_, maxY := g.Size()

	servicesView := views.ViewBox{
		Name: "services",
		Body: `item1
item2
item3`,
		X0: func(i int) int {
			return 1
		},
		Y0: func(i int) int {
			return 1
		},
		X1: func(i int) int {
			return 25
		},
		Y1: func(maxY int) int {
			return maxY / 2
		},
	}

	consoleView := views.ViewBox{
		Name: "console",
		Body: "",
		X0: func(i int) int {
			return 27
		},
		Y0: func(i int) int {
			return 1
		},
		X1: func(maxX int) int {
			return maxX - 1
		},
		Y1: func(maxY int) int {
			return maxY - 1
		},
	}

	//cmdRunner, e := command.NewCommandRunner("docker-compose up", "/home/vmaryn/projects/go/sandbox")
	cmdRunner, e := command.NewCommandRunner("docker exec -it 61d66566deb8 bash", "/home/vmaryn/projects/go/sandbox")

	if e != nil {
		panic(e)
	}
	e = cmdRunner.Start()

	if e != nil {
		panic(e)
	}
	g.SetManager(&servicesView, &consoleView)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		cmdRunner.Stop()
		return gocui.ErrQuit
	}); err != nil {
		log.Panicln(err)
	}

	c := 0
	g.SetKeybinding("", gocui.KeyCtrl2, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		views := []string{"services", "console"}
		g.SetCurrentView(views[c%len(views)])

		c++
		return nil
	})
	g.SetKeybinding("console", gocui.KeyCtrlA, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		view.Autoscroll = !view.Autoscroll
		return nil
	})

	//10
	//150
	if err := g.SetKeybinding("console", gocui.KeyArrowDown, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, y := view.Origin()
		if len(view.BufferLines())-y+1 > maxY {
			view.SetOrigin(x, y+1)
		}

		return nil
	}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("console", gocui.KeyArrowUp, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		x, y := view.Origin()
		if y > 0 {
			view.SetOrigin(x, y-1)
		}

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	g.Update(func(gui *gocui.Gui) error {
		_, e := g.SetCurrentView("services")
		if e != nil {
			log.Panicln(e)
		}

		v, _ := g.View("console")
		v.Editor = gocui.DefaultEditor
		v.Editable = true
		g.SelBgColor = gocui.ColorWhite
		g.SelFgColor = gocui.ColorBlue
		g.Highlight = true
		return nil
	})

	go func() {
		for {
			select {
			case v, ok := <-cmdRunner.OutStream:
				if !ok {
					cmdRunner.Stop()
					return
				}
				consoleView.Body = consoleView.Body + v

				g.Update(func(gui *gocui.Gui) error {
					view, _ := g.View("console")
					view.Clear()
					fmt.Fprintf(view, consoleView.Body)
					return nil
				})
			}
		}

	}()
	go func() {
		for {
			select {
			case v, ok := <-cmdRunner.ErrStream:
				if !ok {
					cmdRunner.Stop()
					return
				}
				consoleView.Body = consoleView.Body + v

				g.Update(func(gui *gocui.Gui) error {
					view, _ := g.View("console")
					view.Clear()
					fmt.Fprintf(view, consoleView.Body)
					return nil
				})
			}
		}

	}()
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
