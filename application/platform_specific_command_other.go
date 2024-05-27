// +build windows darwin

package application

import "os/exec"

func beforeExecuteCommand(cmd *exec.Cmd) {
    // we need to find a way to kill the process when the parent process is killed
    // on windows and macos
}
