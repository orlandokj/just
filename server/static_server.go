package server

import (
	"fmt"
	"github.com/orlandokj/just/config"
	"log"
	"net/http"

	"github.com/urfave/cli/v2"
)


func enableCors(fs http.Handler, config config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" || r.Method == "GET" {
            var origin string
            if config.Cors.DynamicOrigin {
                origin = r.Header.Get("Origin")
            }
            if origin == "" {
                origin = config.Cors.DefaultOrigin
            }
            if origin != "" {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            } 
            w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprintf("%t", config.Cors.AllowCredentials))
            w.Header().Set("Access-Control-Allow-Headers", config.Cors.Headers)
            w.Header().Set("Access-Control-Allow-Methods", config.Cors.Methods)
        }
        if r.Method != "OPTIONS" {
            fs.ServeHTTP(w, r)
        }
    }
}

func ServeStaticFiles(cCtx *cli.Context, config config.Config) error {
        fs := http.FileServer(http.Dir(config.Dir))
        http.Handle("/", enableCors(fs, config))
        log.Printf("Starting to serve files from dir %s on port %d", config.Dir, config.Port)
        err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
        return err
    }
