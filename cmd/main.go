package main

import (
	"cli"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	g, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	projects := []*cli.Project{
		{
			Name:          "cmdapp",
			Dir:           "/home/vmaryn/projects/go/cli",
			Cmd:           "cat services-view.go",
			IsRunning:     false,
			IsHighlighted: true,
		},
		{
			Name:      "dockerapp",
			Dir:       "/home/vmaryn/projects/go/sandbox",
			Cmd:       "docker-compose up",
			IsRunning: false,
		},
		{
			Name:      "dockerap2p",
			Dir:       "/home/vmaryn/projects/go/sandbox",
			Cmd:       "docker-compose up",
			IsRunning: false,
		},
	}
	app := &cli.AppContext{
		Projects: projects,
		ConBuff: &cli.ConsoleBuff{
			Data: []string{},
		},
	}

	consoleView := cli.NewConsoleView("")
	servicesView := cli.NewServicesView(app.Projects, `item1`)
	app.Console = consoleView
	app.Services = servicesView

	g.SetManager(servicesView, consoleView)

	if e := cli.SetupConsoleBindings(g); e != nil {
		panic(e)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		//cmdRunner.Stop()
		return gocui.ErrQuit
	}); err != nil {
		log.Panicln(err)
	}

	c := 0
	g.SetKeybinding("", gocui.KeyCtrl2, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		views := []string{cli.SERVICES_VIEW, cli.CONSOLE_VIEW}
		g.SetCurrentView(views[c%len(views)])

		c++
		return nil
	})

	g.SetKeybinding(cli.SERVICES_VIEW, gocui.KeyArrowDown, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for i, p := range projects {
			if p.IsHighlighted && i < len(projects)-1 {
				g.Update(func(gui *gocui.Gui) error {
					p.IsHighlighted = false
					projects[i+1].IsHighlighted = true
					servicesView.ReDraw(view)
					return nil
				})
				break
			}
		}

		return nil
	})

	g.SetKeybinding(cli.SERVICES_VIEW, gocui.KeyArrowUp, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for i, p := range projects {
			if p.IsHighlighted && i != 0 {
				g.Update(func(gui *gocui.Gui) error {
					p.IsHighlighted = false
					projects[i-1].IsHighlighted = true
					servicesView.ReDraw(view)
					return nil
				})

				break
			}
		}

		return nil
	})
	g.SetKeybinding(cli.SERVICES_VIEW, gocui.KeyCtrlR, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		for _, p := range projects {
			if p.IsHighlighted {
				if p.IsRunning {
					p.CmdInst.Stop()
				} else {
					e := p.Start()
					if e != nil {
						log.Panicln(e)
					}
					g.Update(func(gui *gocui.Gui) error {
						servicesView.ReDraw(view)

						return nil
					})
					go func() {
						v, _ := g.View(cli.CONSOLE_VIEW)

						g.Update(func(gui *gocui.Gui) error {
							v.Clear()
							fmt.Fprintf(v, p.Data)
							return nil
						})

						for {
							_, ok := <-p.DataChanged

							if !ok {
								g.Update(func(gui *gocui.Gui) error {
									servicesView.ReDraw(view)

									return nil
								})
								return
							}
							g.Update(func(gui *gocui.Gui) error {
								v.Clear()
								fmt.Fprintf(v, p.Data)
								return nil
							})
						}

					}()
				}

				break
			}
		}
		return nil
	})

	g.Update(func(gui *gocui.Gui) error {
		_, e := g.SetCurrentView(cli.SERVICES_VIEW)
		if e != nil {
			log.Panicln(e)
		}

		v, _ := g.View(cli.CONSOLE_VIEW)
		v.Editor = gocui.DefaultEditor
		v.Editable = true
		g.SelBgColor = gocui.ColorWhite
		g.SelFgColor = gocui.ColorBlue
		g.Highlight = true
		return nil
	})

	//go func() {
	//	for {
	//		select {
	//		case v, ok := <-cmdRunner.OutStream:
	//			if !ok {
	//				cmdRunner.Stop()
	//				return
	//			}
	//			consoleView.Body = consoleView.Body + v
	//
	//			g.Update(func(gui *gocui.Gui) error {
	//				view, _ := g.View("console")
	//				view.Clear()
	//				fmt.Fprintf(view, consoleView.Body)
	//				return nil
	//			})
	//		}
	//	}
	//
	//}()
	//go func() {
	//	for {
	//		select {
	//		case v, ok := <-cmdRunner.ErrStream:
	//			if !ok {
	//				cmdRunner.Stop()
	//				return
	//			}
	//			consoleView.Body = consoleView.Body + v
	//
	//			g.Update(func(gui *gocui.Gui) error {
	//				view, _ := g.View("console")
	//				view.Clear()
	//				fmt.Fprintf(view, consoleView.Body)
	//				return nil
	//			})
	//		}
	//	}
	//
	//}()
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

//
//func StartProject(app *cli.AppContext, g *gocui.Gui, project *cli.Project) {
//	active := app.Active
//
//	if active != nil {
//		active.IsRunning = false
//		active.CmdInst.Stop()
//	}
//
//	cmd, _ := command.NewCommandRunner(project.Cmd, project.Dir)
//	project.CmdInst = *cmd
//	project.CmdInst.Start()
//	project.IsRunning = true
//	app.Active = project
//
//	go func() {
//		for {
//			select {
//			case v, ok := <-project.CmdInst.OutStream:
//				if !ok {
//					return
//				}
//				app.Console.Body = app.Console.Body + v
//
//				g.Update(func(gui *gocui.Gui) error {
//					view, _ := g.View(cli.CONSOLE_VIEW)
//					view.Clear()
//					fmt.Fprintf(view, app.Console.Body)
//					return nil
//				})
//			case <-project.CmdInst.Done:
//				return
//			}
//		}
//
//	}()
//}

func layout(g *gocui.Gui) error {

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
