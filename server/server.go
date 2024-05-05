package server

import (
	"errors"

	"github.com/orlandokj/just/config"
)

type Server interface {
    Build() error
    Run() error
}


func CreateServer(config config.Config) (Server, error) {
    switch config.Type {
    case "java":
        return CreateJavaServer(config)
    case "static":
        return CreateStaticServer(config)
    default:
        return nil, errors.New("Unknown server type: " + config.Type)
    }
}
