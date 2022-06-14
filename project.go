package cli

import (
	"cli/command"
)

type Project struct {
	Name          string
	Dir           string
	Cmd           string
	IsRunning     bool
	IsHighlighted bool
	CmdInst       command.CommandRunner
	Data          string
	DataChanged   chan struct{}
}

func (p *Project) Start() error {
	cmd, _ := command.NewCommandRunner(p.Cmd, p.Dir)
	p.CmdInst = *cmd
	p.DataChanged = make(chan struct{})

	if e := p.CmdInst.Start(); e != nil {
		return e
	}

	p.IsRunning = true

	go func() {
		for {
			select {
			case v, ok := <-p.CmdInst.OutStream:
				if !ok {
					p.IsRunning = false
					return
				}

				p.Data += v
				p.DataChanged <- struct{}{}
			case <-p.CmdInst.Done:
				p.IsRunning = false
				return
			}
		}

	}()

	return nil
}
