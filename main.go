package main

import (
	"log"
	"os"

	"github.com/orlandokj/just/config"
	"github.com/orlandokj/just/server"
	"github.com/urfave/cli/v2"
)

func getServer(config config.Config) (server.Server, error) {
    return server.CreateServer(config)
}

func runServer(config config.Config) error {
    s, err := getServer(config)
    if err != nil {
        return err
    }
    return s.Run()
}

func buildServer(config config.Config) error {
    s, err := getServer(config)
    if err != nil {
        return err
    }
    return s.Build()
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
                    return runServer(loadedConfig)
                },
            },
            {
                Name:    "build",
                Usage:   "Build the project to run after",
                Action: func(cCtx *cli.Context) error {
                    return buildServer(loadedConfig)
                },
            },
            {
                Name:    "build-run",
                Usage:   "Build the project and run after the build completed",
                Action: func(cCtx *cli.Context) error {
                    server, err := getServer(loadedConfig)
                    if err != nil {
                        return err
                    }
                    err = server.Build()
                    if err != nil {
                        return err
                    }
                    return server.Run()
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
            return runServer(loadedConfig)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
