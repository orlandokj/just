package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/orlandokj/just/config"
	"github.com/orlandokj/just/server"
	"github.com/urfave/cli/v2"
)

func runServer(cCtx *cli.Context, config config.Config) error {
    if config.Type == "static" {
        return server.ServeStaticFiles(cCtx, config)
    }

    return errors.New(fmt.Sprintf("Invalid server type: %s", config.Type))
}

func main() {
	var loadedConfig config.Config
    var configFile string
	app := &cli.App{
		Name:  "Just CLI",
		Usage: "A CLI application to just run your projects",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   "just.config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
                Destination: &configFile,
			},
		},
        Commands: []*cli.Command{
            {
                Name:    "run",
                Usage:   "Run the server",
                Action: func(cCtx *cli.Context) error {
                    return runServer(cCtx, loadedConfig)
                },
            },
            {
                Name:    "build",
                Usage:   "Build the project to run after",
                Action: func(cCtx *cli.Context) error {
                    log.Println("Building project")
                    return nil
                },
            },
        },
        Before: func(cCtx *cli.Context) error {
            var err error
            log.Printf("Loading config from: %s\n", configFile)
            loadedConfig, err = config.LoadConfig(configFile)
            return err

        },
		Action: func(cCtx *cli.Context) error {
            return runServer(cCtx, loadedConfig)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
