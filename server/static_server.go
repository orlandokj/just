package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/orlandokj/just/config"
)


func enableCors(fs http.Handler, config CorsConfig) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" || r.Method == "GET" {
            var origin string
            if config.DynamicOrigin {
                origin = r.Header.Get("Origin")
            }
            if origin == "" {
                origin = config.DefaultOrigin
            }
            if origin != "" {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            } 
            w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprintf("%t", config.AllowCredentials))
            w.Header().Set("Access-Control-Allow-Headers", config.Headers)
            w.Header().Set("Access-Control-Allow-Methods", config.Methods)
        }
        if r.Method != "OPTIONS" {
            fs.ServeHTTP(w, r)
        }
    }
}

type CorsConfig struct {
    DynamicOrigin bool `json:"dynamicOrigin"`
    DefaultOrigin string `json:"defaultOrigin"`
    Headers string `json:"headers"`
    Methods string `json:"methods"`
    AllowCredentials bool `json:"allowCredentials"`
}

type StaticServerConfig struct {
    Dir string `json:"dir"`
    Port int `json:"port"`
    Cors CorsConfig `json:"cors"`
}

type StaticServer struct {
    config StaticServerConfig
    logFunc LogFunc
}

func (ss StaticServer) Build() (ServerProcess, error) {
    return nil, errors.New("Static server does not support build yet")
}

func (ss StaticServer) Run() (ServerProcess, error) {
    serveDir := ss.config.Dir
    if serveDir == "" {
        serveDir = "./dist"
    }
    fs := http.FileServer(http.Dir(serveDir))
    if ss.config.Cors == (CorsConfig{}) {
       ss.config.Cors = CorsConfig{
            DynamicOrigin: true,
            DefaultOrigin: "http://localhost:9000",
            Headers: "x-dispositivo, x-requested-with",
            Methods: "*",
            AllowCredentials: true,
        }
    }
    mux := http.NewServeMux()
    mux.Handle("/", enableCors(fs, ss.config.Cors))
    log.Printf("Starting to serve files from dir %s on port %d", ss.config.Dir, ss.config.Port)
    runningServer := RunningServer{
        server: &http.Server{
            Addr: fmt.Sprintf(":%d", ss.config.Port),
            Handler: mux,
        },
    }
    go func() {
        runningServer.server.ListenAndServe()
    }()
    return runningServer, nil
}

type RunningServer struct {
    server *http.Server
    stopped bool
}

func (m RunningServer) Stop() error {
    if m.stopped == true {
        return nil
    }

    err := m.server.Close()
    m.stopped = true
    return err
}

func (m RunningServer) MemoryUsage() int {
    return -1
}

func (m RunningServer) CPUUsage() int {
    return -1
}

func CreateStaticServer(config config.Config, logFunc LogFunc) (Server, error) {
    staticConfig := StaticServerConfig{}
    err := config.ToConfigType(&staticConfig)
    if err != nil {
        return nil, err
    }
    
    return StaticServer{
        config: staticConfig,
        logFunc: logFunc,
    }, nil
}
