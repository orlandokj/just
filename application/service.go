package application

import (
	"container/list"
	"errors"
	"log"
	"sort"
)

type Application struct {
    Name string `json:"name"`
    Config ApplicationConfig `json:"config"`
    process RunningProcess
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


    applicationHandler, err := CreateApplication(a.Config, func(log string) {
        a.logs.PushBack(log)
        if a.logChan != nil {
            a.logChan <- log
        }
    })
    if err != nil {
        return err
    }

    a.Status = "Running"
    process, err := applicationHandler.Run()
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
    a.logs.Init()

    return nil
}

func NewApplication(config ApplicationConfig) error {
    if config.Name == "" {
        return errors.New("Name is required")
    }

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

func ModifyApplication(config ApplicationConfig) error {
    app, ok := applications[config.Name]
    if !ok {
        return errors.New("Application not found")
    }

    app.Config = config
    return EditApplication(*app)
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

func GetApplications() []*Application {
    var apps []*Application
    for _, app := range applications {
        apps = append(apps, app)
    }
    sort.Slice(apps, func(i, j int) bool {
        return apps[i].Name < apps[j].Name
    })
    return apps
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

func GetApplication(name string) *Application {
    app, ok := applications[name]
    if !ok {
        return nil
    }

    return app
}

func init() {
    loadApplications()
}
