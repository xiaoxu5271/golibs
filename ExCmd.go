package golibs

import (
	"context"
	"os"
	"os/exec"
)

type ExCmd struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
	state  chan string

	er CmdBuf
	ou CmdBuf
}

func (e *ExCmd) Init(name string, arg ...string) *ExCmd {
	ctx, cancel := context.WithCancel(context.Background())
	e.cmd = exec.CommandContext(ctx, name, arg...)
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

	e.cmd.Stdout = &e.ou
	e.cmd.Stderr = &e.er

	e.cancel = cancel
	e.state = make(chan string, 1)
	return e
}

func (e *ExCmd) InitPath(dir, name string, arg ...string) *ExCmd {
	e.Init(name, arg...)
	e.cmd.Dir = dir
	return e
}

func (e *ExCmd) exec() {
	err := e.cmd.Run()
	if err != nil {
		e.state <- "failed"
	} else {
		e.state <- "finished"
	}
}

func (e *ExCmd) Exec() {
	select {
	case <-e.state:
	default:
	}
	go e.exec()
}

func (e *ExCmd) GetWaitEvent() <-chan string {
	return e.state
}

func (e *ExCmd) Output() string {
	return e.ou.TruncateString()
}

func (e *ExCmd) Error() string {
	return e.er.TruncateString()
}

func (e *ExCmd) ClearOutput() {
	e.ou.Reset()
}

func (e *ExCmd) ClearError() {
	e.er.Reset()
}

func (e *ExCmd) Kill() {
	if e.cancel != nil {
		e.cancel()
	}
}

func (e *ExCmd) WaitDone() bool {
	ste := <-e.state
	return ste == "finished"
}
