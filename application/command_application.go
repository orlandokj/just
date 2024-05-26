package application

import (
	"bufio"
	"os"
	"os/exec"
)

type RunningCommand struct {
    pid int
    stopped bool
}

func (rc RunningCommand) Stop() error {
    if rc.stopped == true {
        return nil
    }

    p, err := os.FindProcess(rc.pid)
    if err != nil {
        return err
    }
    return p.Kill()
}

func (rc RunningCommand) MemoryUsage() int {
    // TODO How to get memory usage of a process?
    return -1
}

func (rc RunningCommand) CPUUsage() int {
    // TODO How to get CPU usage of a process?
    return -1
}

func RunCmd(cmd *exec.Cmd, logFunc LogFunc) (RunningProcess, error) {
    beforeExecuteCommand(cmd)
    stdout, err := cmd.StdoutPipe()
    cmd.Stderr = cmd.Stdout
    if err != nil {
        return nil, err
    }

    err = cmd.Start()
    if err != nil {
        logFunc(err.Error())
        return nil, err
    }

    go func() {
        reader := bufio.NewReader(stdout)

        for {
            line, err := reader.ReadString('\n')
            if err != nil {
                break
            }
            logFunc(line)

        }
        cmd.Wait()
    }()
        
    return &RunningCommand{
        pid: cmd.Process.Pid,
    }, nil
}

