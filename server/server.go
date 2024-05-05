package server

import (
	"errors"

	"github.com/orlandokj/just/config"
)

type ServerProcess interface {
    Stop() error
    MemoryUsage() int
    CPUUsage() int
}

type Server interface {
    Build() (ServerProcess, error)
    Run() (ServerProcess, error)
}

type LogFunc func(string)

func CreateServer(config config.Config, logFunc LogFunc) (Server, error) {
    switch config.Type {
    case "java":
        return CreateJavaServer(config, logFunc)
    case "static":
        return CreateStaticServer(config, logFunc)
    default:
        return nil, errors.New("Unknown server type: " + config.Type)
    }
}
