package golibs

import (
	"os"
	"os/exec"
)

type SimCmd struct {
	cmd   *exec.Cmd
	state chan string
}

func (e *SimCmd) Init(name string, arg ...string) *SimCmd {
	e.cmd = exec.Command(name, arg...)
	envs := []string{
		"USER",
		"HOME",
		"SHELL",
		"SHLVL",
		"LOGNAME",
		"PATH",
	}
	for _, key := range envs {
		env := os.Getenv(key)
		e.cmd.Env = append(e.cmd.Env, key+"="+env)
	}

	e.cmd.Stdout = nil
	e.cmd.Stderr = nil

	e.state = make(chan string, 1)
	return e
}

func (e *SimCmd) InitPath(dir, name string, arg ...string) *SimCmd {
	e.Init(name, arg...)
	e.cmd.Dir = dir
	return e
}

func (e *SimCmd) exec() {
	err := e.cmd.Run()
	if err != nil {
		e.state <- "failed"
	} else {
		e.state <- "finished"
	}
}

func (e *SimCmd) Exec() {
	select {
	case <-e.state:
	default:
	}
	go e.exec()
}

func (e *SimCmd) GetWaitEvent() <-chan string {
	return e.state
}

func (e *SimCmd) WaitDone() bool {
	ste := <-e.state
	return ste == "finished"
}
