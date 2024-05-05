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
}

func (ss StaticServer) Build() error {
    return errors.New("Static server does not support build yet")
}

func (ss StaticServer) Run() error {
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
    http.Handle("/", enableCors(fs, ss.config.Cors))
    log.Printf("Starting to serve files from dir %s on port %d", ss.config.Dir, ss.config.Port)
    err := http.ListenAndServe(fmt.Sprintf(":%d", ss.config.Port), nil)
    return err
}

func CreateStaticServer(config config.Config) (Server, error) {
    staticConfig := StaticServerConfig{}
    err := config.ToConfigType(&staticConfig)
    if err != nil {
        return nil, err
    }
    
    return StaticServer{
        config: staticConfig,
    }, nil
}
