package application

import (
	"errors"

)

type RunningProcess interface {
    Stop() error
    MemoryUsage() int
    CPUUsage() int
}

type ApplicationHandler interface {
    Build() (RunningProcess, error)
    Run() (RunningProcess, error)
}

type LogFunc func(string)

func CreateApplication(config ApplicationConfig, logFunc LogFunc) (ApplicationHandler, error) {
    switch config.Type {
    case "java":
        return CreateJavaApplication(config, logFunc)
    case "static":
        return CreateStaticApplication(config, logFunc)
    default:
        return nil, errors.New("Unknown application type: " + config.Type)
    }
}
