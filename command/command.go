package command

import (
	"fmt"
	"github.com/creack/pty"
	"io"
	"os"
	"os/exec"
	"strings"
)

type CommandRunner struct {
	Cmd       *exec.Cmd
	OutStream chan string
	ErrStream chan string
	stdOut    io.ReadCloser
	stdErr    io.ReadCloser
	P         *os.File
}

func NewCommandRunner(bashCmd string, dir string) (*CommandRunner, error) {
	var cmd *exec.Cmd
	var runner *CommandRunner

	//if !strings.Contains(bashCmd, "docker") {
	//	cmd = exec.Command("/bin/bash", "-c", bashCmd)
	//	stdOut, _ := cmd.StdoutPipe()
	//	stdErr, _ := cmd.StderrPipe()
	//	runner = &CommandRunner{
	//		Cmd:       cmd,
	//		OutStream: make(chan string, 10),
	//		ErrStream: make(chan string, 10),
	//		stdOut:    stdOut,
	//		stdErr:    stdErr,
	//	}
	//	cmd.Dir = dir
	//
	//	panic(bashCmd)
	//} else {
	args := strings.Split(bashCmd, " ")
	cmd = exec.Command(args[0], args[1:]...)
	//b := new(bytes.Buffer)
	//cmd.Stdin = b
	p, e := pty.Start(cmd)

	if e != nil {
		fmt.Println(e)
	}

	runner = &CommandRunner{
		Cmd:       cmd,
		OutStream: make(chan string, 10),
		ErrStream: make(chan string, 10),
		P:         p,
		stdOut:    p,
		stdErr:    p,
	}
	//}

	//stdOut, e := cmd.StdoutPipe()
	//stdErr, e := cmd.StderrPipe()

	//if e != nil {
	//	return nil, e
	//}

	return runner, nil
}

func (c *CommandRunner) Stop() error {
	return c.Cmd.Process.Kill()
}

func (c *CommandRunner) Start() error {
	//e := c.Cmd.Start()
	//
	//if e != nil {
	//	return e
	//}

	go func() {
		b := make([]byte, 1024)
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
	//
	//go func() {
	//	b := make([]byte, 64)
	//	read := 0
	//	var e error
	//
	//	for {
	//		read, e = c.stdErr.Read(b)
	//		if e != nil {
	//			if e == io.EOF {
	//				close(c.ErrStream)
	//				return
	//			}
	//
	//			panic(e)
	//		} else {
	//			c.ErrStream <- string(b[:read])
	//		}
	//	}
	//}()

	return nil
}
