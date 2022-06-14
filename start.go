package cli

import (
	"cli/command"
)

func StartProject(project *Project) {
	cmd, _ := command.NewCommandRunner(project.Cmd, project.Dir)
	project.CmdInst = *cmd
	project.CmdInst.Start()
	project.IsRunning = true

	go func() {
		for {
			select {
			case v, ok := <-project.CmdInst.OutStream:
				if !ok {
					return
				}
				project.Data += v
				project.DataChanged <- struct{}{}
			case <-project.CmdInst.Done:
				project.IsRunning = false
				return
			}
		}

	}()
}
