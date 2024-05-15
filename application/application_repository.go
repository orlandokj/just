package application

import (
	"encoding/json"
	"log"
	"os"

)

var configDir = ".config/just"
var applicationFile = "applications.json"

var applicationCache []Application

func getApplications() []Application {
    if applicationCache == nil {
        applicationCache = loadApplicationsFile()
    }

    return applicationCache
}

func AddApplication(application Application) error {
    applications := getApplications()
    applications = append(applications, application)
    err := saveApplicationsFile(applications)
    if err != nil {
        return err
    }

    return nil
}

func RemoveApplication(application Application) error {
    applications := getApplications()
    for i, app := range applications {
        if app.Name == application.Name {
            applications = append(applications[:i], applications[i+1:]...)
            break
        }
    }

    err := saveApplicationsFile(applications)
    if err != nil {
        return err
    }

    return nil
}

func loadApplicationsFile() []Application {
    applications := make([]Application, 0)
    userDir, err := os.UserHomeDir()
    if err != nil {
        log.Printf("Error getting user home dir: %v", err)
        return applications
    }

    err = os.MkdirAll(userDir + "/" + configDir, 0664)
    if err != nil {
        log.Printf("Error creating config dir: %v", err)
        return applications
    }
    configFile, err := os.OpenFile(userDir + "/" + configDir + "/" + applicationFile, os.O_CREATE|os.O_RDONLY, 0664)
    if err != nil {
        log.Printf("Error opening applications file: %v", err)
        return applications
    }

    err = json.NewDecoder(configFile).Decode(&applications)
    if err != nil {
        log.Printf("Error parsing json applications file: %v", err)
    }

    return applications
}

func saveApplicationsFile(applications []Application) error {
    userDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    fileContent, err := json.Marshal(applications)
    if err != nil {
        return err
    }

    err = os.WriteFile(userDir + "/" + configDir + "/" + applicationFile, fileContent, 0644)
    if err != nil {
        return err
    }

    return nil
}
