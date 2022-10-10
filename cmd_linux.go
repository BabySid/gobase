package gobase

import (
	"os/exec"
	"syscall"
)

func SetChildrenProcessDetached(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{}
}

func SetChildrenProcessGroupID(c *exec.Cmd) {
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.Setpgid = true
	c.SysProcAttr.Pdeathsig = syscall.SIGTERM
}
