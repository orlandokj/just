// +build !windows

package application

import (
	"os/exec"
	"syscall"
)

func beforeExecuteCommand(cmd *exec.Cmd) {
    cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}
}
