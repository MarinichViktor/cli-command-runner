package cli

import (
	"cli/command"
	"strings"
	"time"
)

const BUFF_LIMIT = 512

type ProjectArgs struct {
	Name string `yaml:"name"`
	Dir  string `yaml:"dir"`
	Cmd  string `yaml:"cmd"`
}

type Project struct {
	Name               string
	Dir                string
	Cmd                string
	IsRunning          bool
	IsHighlighted      bool
	CmdInst            command.CommandRunner
	Data               []string
	DataChanged        chan struct{}
	Subscriptions      map[uint]*Subscription
	lastSubscriptionId uint
	ViewName           string
	HasSubscription    bool
	LastSynced         time.Time
	LastUpdated        time.Time
	View               *VirtualView
}

func NewProject(a *ProjectArgs) *Project {
	return &Project{
		Name:          a.Name,
		Dir:           a.Dir,
		Cmd:           a.Cmd,
		IsRunning:     false,
		IsHighlighted: false,
		Subscriptions: make(map[uint]*Subscription),
		LastUpdated:   time.Now(),
		LastSynced:    time.Now(),
		View: &VirtualView{
			Offset: 0,
			Data:   []string{},
		},
	}
}

type Subscription struct {
	Data func(string)
	Done func()
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
	p.View = &VirtualView{
		Offset: 0,
		Data:   []string{},
	}

	if e := p.CmdInst.Start(); e != nil {
		return e
	}

	p.IsRunning = true
	done := make(chan struct{})

	go func() {
		for {
			select {
			// todo remove select
			case v, ok := <-p.CmdInst.OutStream:
				if !ok {
					for _, s := range p.Subscriptions {
						s.Done()
					}

					p.IsRunning = false
					done <- struct{}{}
					return
				}

				newData := strings.Split(v, "\n")

				//for i, s := range newData {
				//	newData[i] = strings.Trim(s, "\x00")
				//}

				l := len(p.View.Data)

				if l > 0 {
					p.View.Data[l-1] += newData[0]
					p.View.AppendData(newData[1:])
				} else {
					p.View.AppendData(newData)
				}
				//p.Data = append(p.Data, newData[1:]...)

				//if l > BUFF_LIMIT {
				//	p.Data = p.Data[l-BUFF_LIMIT : l-1]
				//}
				p.LastUpdated = time.Now()
				//for _, s := range p.Subscriptions {
				//	s.Data(p.StrData())
				//}
			}
		}

	}()
	go func() {
		for {
			t := time.NewTicker(150 * time.Millisecond)
			select {
			case <-t.C:
				if p.LastSynced.Before(p.LastUpdated) {
					for _, s := range p.Subscriptions {
						s.Data(p.StrData())
					}
					p.LastSynced = time.Now()
				}
			case <-done:
				return
			}
		}
	}()

	return nil
}

func (p *Project) StrData() string {
	return strings.Join(p.Data, "\n")
}

func (p *Project) Stop() error {
	if e := p.CmdInst.Stop(); e != nil {
		return e
	}

	return nil
}
