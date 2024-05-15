package application

import (
	"container/list"
	"errors"
	"log"

	"github.com/orlandokj/just/config"
	"github.com/orlandokj/just/server"
)

type Application struct {
    Name string `json:"name"`
    Config config.Config `json:"config"`
    process server.ServerProcess
    Status string
    logs list.List
    logChan chan string
}

func (a *Application) MemoryUsage() int {
    if a.process == nil {
        return -1
    }
    return a.process.MemoryUsage()
}

func (a *Application) CPUUsage() int {
    if a.process == nil {
        return -1
    }

    return a.process.CPUUsage()
}

var applications = make(map[string]*Application)

func (a *Application) run() error {
    if a.Status == "Running" {
        return errors.New("Application already running")
    }


    server, err := server.CreateServer(a.Config, func(log string) {
        a.logs.PushBack(log)
        if a.logChan != nil {
            a.logChan <- log
        }
    })
    if err != nil {
        return err
    }

    a.Status = "Running"
    process, err := server.Run()
    if err != nil {
        a.Status = "Failed to run"
        return err
    }

    a.process = process
    return nil
}

func (a *Application) stop() error {
    if a.Status != "Running" {
        return errors.New("Application not running")
    }

    err := a.process.Stop()

    if err != nil {
        return err
    }

    a.Status = "Stopped"

    return nil
}

func NewApplication(config config.Config) error {
    app := Application{
        Name: config.Name,
        Config: config,
        Status: "Stopped",
    }

    err := AddApplication(app)
    if err != nil {
        return err
    }
    applications[config.Name] = &app
    return nil
}

func RunApplication(name string) error {
    app, ok := applications[name]
    if !ok {
        return errors.New("Application not found")
    }

    return app.run()
}

func StopApplication(name string) error {
    app, ok := applications[name]
    if !ok {
        return errors.New("Application not found")
    }

    return app.stop()
}

func loadApplications() {
    log.Println("Loading applications")
    apps := getApplications()
    for _, app := range apps {
        app.Status = "Stopped"
        applications[app.Name] = &app
    } 
}

func GetApplications() map[string]*Application {
    return applications
}

func DeleteApplication(name string) error {
    app, ok := applications[name]
    if !ok {
        return nil
    }

    if app.Status == "Running" {
        return errors.New("Application is running")
    }

    delete(applications, name)
    return RemoveApplication(*app)
}

func WatchLogs(name string, logChan chan string) error {
    app, ok := applications[name]
    if !ok || app.Status == "Stopped" {
        return errors.New("Application not found")
    }

    go func() {
        for e := app.logs.Front(); e != nil; e = e.Next() {
            logChan <- e.Value.(string)
        }
    }()
    app.logChan = logChan

    return nil
}

func StopWatching(name string) {
    app, ok := applications[name]
    if !ok {
        return
    }

    app.logChan = nil
}

func init() {
    loadApplications()
}
