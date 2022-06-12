package command

import (
	"io"
	"os/exec"
)

type CommandRunner struct {
	Cmd       *exec.Cmd
	OutStream chan string
	ErrStream chan string
	stdOut    io.ReadCloser
	stdErr    io.ReadCloser
}

func NewCommandRunner(bashCmd string, dir string) (*CommandRunner, error) {
	cmd := exec.Command("/bin/bash", "-c", bashCmd)
	cmd.Dir = dir
	stdOut, e := cmd.StdoutPipe()
	stdErr, e := cmd.StderrPipe()

	if e != nil {
		return nil, e
	}

	runner := &CommandRunner{
		Cmd:       cmd,
		OutStream: make(chan string, 10),
		ErrStream: make(chan string, 10),
		stdOut:    stdOut,
		stdErr:    stdErr,
	}

	return runner, nil
}

func (c *CommandRunner) Stop() error {
	return c.Cmd.Process.Kill()
}

func (c *CommandRunner) Start() error {
	e := c.Cmd.Start()

	if e != nil {
		return e
	}

	go func() {
		b := make([]byte, 128)
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
					close(c.ErrStream)
					return
				}

				panic(e)
			} else {
				c.ErrStream <- string(b[:read])
			}
		}
	}()

	return nil
}
