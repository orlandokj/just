package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

type Config struct {
    Name    string `json:"name"`
    Port    int    `json:"port"`
    Type    string `json:"type"`
    Dir     string `json:"dir"`
}

func loadConfig(configPath string) (Config, error) {
    filePath := configPath
    if filePath == "" {
        filePath = "just.config"
    }

    log.Println("Carregando configurações de", filePath)

    file, err := os.Open(filePath)
    config := Config{
        Dir: "./dist",
    }

    if err != nil {
        return config, err
    }
    defer file.Close()

    json.NewDecoder(file).Decode(&config)
    return config, nil
}

func enableCors(fs http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" || r.Method == "GET" {
            origin := r.Header.Get("Origin")
            if origin == "" {
                origin = "http://localhost:9000"
            }
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Allow-Headers", "x-dispositivo, x-requested-with")
        }
        if r.Method != "OPTIONS" {
            fs.ServeHTTP(w, r)
        }
    }
}

func main() {
    app := &cli.App{
        Name:  "Just CLI",
        Usage: "A CLI application to just run your projects",
        Action: func(cCtx *cli.Context) error {
            config, err := loadConfig(cCtx.String("config"))
            if err != nil {
                return err
            }

            fs := http.FileServer(http.Dir(config.Dir))
            http.Handle("/", enableCors(fs))
            log.Printf("Iniciando aplicação %s na porta %d", config.Name, config.Port)
            err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
            return err
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
