package main

import (
	"os"
	"os/exec"
	"syscall"
)

func newCmd(name string, arg ...string) *Cmd {
	c := new(Cmd)
	c.init(name, arg...)
	return c
}

type Cmd struct {
	*exec.Cmd
}

func (c *Cmd) Stop() error {
	return syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
}

func (c *Cmd) Restart() error {
	c.Stop()
	c.Wait()
	c.init(c.Args[0], c.Args[1:]...)
	return c.Start()
}

func (c *Cmd) init(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	c.Cmd = cmd
}
