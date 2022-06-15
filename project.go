package cli

import (
	"cli/command"
)

type ProjectArgs struct {
	Name string `yaml:"name"`
	Dir  string `yaml:"dir"`
	Cmd  string `yaml:"cmd"`
}

func NewProject(pArgs *ProjectArgs) *Project {
	return &Project{
		Name:          pArgs.Name,
		Dir:           pArgs.Dir,
		Cmd:           pArgs.Cmd,
		IsRunning:     false,
		IsHighlighted: false,
	}
}

type Subscription struct {
	Data func(string)
	Done func()
}

type Project struct {
	Name               string
	Dir                string
	Cmd                string
	IsRunning          bool
	IsHighlighted      bool
	CmdInst            command.CommandRunner
	Data               string
	DataChanged        chan struct{}
	Subscriptions      map[uint]*Subscription
	lastSubscriptionId uint
}

func (p *Project) Subscribe(s func(string), d func()) func() {
	p.lastSubscriptionId++
	p.Subscriptions[p.lastSubscriptionId] = &Subscription{
		Data: s,
		Done: d,
	}

	return func() {
		delete(p.Subscriptions, p.lastSubscriptionId)
	}
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
				for _, s := range p.Subscriptions {
					s.Data(p.Data)
				}
			case <-p.CmdInst.Done:
				p.IsRunning = false
				for _, s := range p.Subscriptions {
					s.Done()
				}
				return
			}
		}

	}()

	return nil
}
