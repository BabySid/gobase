package gobase

import (
	"os/exec"
	"syscall"
)

func SetChildrenProcessDetached(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{}
}

func SetChildrenProcessGroupID(c *exec.Cmd) {
}
