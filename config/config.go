package config

import (
	"encoding/json"
	"os"
)

type Config struct {
    Type    string `json:"type"`
    Name    string `json:"name"`
    WorkDir string `json:"workDir"`
    Config  map[string]interface{} `json:"config"`
}

func (c Config) ToConfigType(output interface{}) error  {
    jsonBody, err := json.Marshal(c.Config)
    if err != nil {
        return err
    }
    return json.Unmarshal(jsonBody, output)
}

func LoadConfig(configPath string) (Config, error) {
    filePath := configPath
    if filePath == "" {
        filePath = "just.config"
    }

    file, err := os.Open(filePath)
    config := Config{}

    if err != nil {
        return config, err
    }
    defer file.Close()

    json.NewDecoder(file).Decode(&config)
    return config, nil
}

