package command

import (
	"io"
	"os/exec"
)

type CommandRunner struct {
	Cmd       *exec.Cmd
	OutStream chan string
	ErrStream chan string
	Done      chan struct{}
	stdOut    io.ReadCloser
	stdErr    io.ReadCloser
}

func NewCommandRunner(bashCmd string, dir string) (*CommandRunner, error) {
	var cmd *exec.Cmd
	var runner *CommandRunner
	cmd = exec.Command("/bin/bash", "-c", bashCmd)

	stdOut, e := cmd.StdoutPipe()
	if e != nil {
		return nil, e
	}

	stdErr, e := cmd.StderrPipe()
	if e != nil {
		return nil, e
	}

	runner = &CommandRunner{
		Cmd:       cmd,
		OutStream: make(chan string, 16),
		ErrStream: make(chan string, 16),
		stdOut:    stdOut,
		stdErr:    stdErr,
		Done:      make(chan struct{}),
	}
	cmd.Dir = dir

	return runner, nil
}

func (c *CommandRunner) Stop() error {
	return c.Cmd.Process.Kill()
}

func (c *CommandRunner) Start() error {
	if e := c.Cmd.Start(); e != nil {
		return e
	}

	go func() {
		c.Cmd.Wait()
		c.Done <- struct{}{}
	}()

	go func() {
		b := make([]byte, 256)
		read := 0
		var e error

		for {
			read, e = c.stdOut.Read(b)
			if e != nil {
				if e == io.EOF {
					close(c.OutStream)
					return
				}

				panic(e)
			} else {
				c.OutStream <- string(b[:read])
			}
		}
	}()

	go func() {
		b := make([]byte, 64)
		read := 0
		var e error

		for {
			read, e = c.stdErr.Read(b)
			if e != nil {
				if e == io.EOF {
					return
				}

				panic(e)
			} else {
				c.OutStream <- string(b[:read])
			}
		}
	}()

	return nil
}
