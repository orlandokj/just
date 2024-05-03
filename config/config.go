package config

import (
	"encoding/json"
	"os"
)

type CorsConfig struct {
    DynamicOrigin bool `json:"dynamicOrigin"`
    DefaultOrigin string `json:"defaultOrigin"`
    Headers string `json:"headers"`
    Methods string `json:"methods"`
    AllowCredentials bool `json:"allowCredentials"`
}

type Config struct {
    Name    string `json:"name"`
    Port    int    `json:"port"`
    Type    string `json:"type"`
    Dir     string `json:"dir"`
    Cors    CorsConfig `json:"cors"`
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
    config := Config{
        Dir: "./dist",
        Cors: CorsConfig{
            DynamicOrigin: true,
            DefaultOrigin: "http://localhost:9000",
            Headers: "x-dispositivo, x-requested-with",
            Methods: "*",
            AllowCredentials: true,
        },

    }

    if err != nil {
        return config, err
    }
    defer file.Close()

    json.NewDecoder(file).Decode(&config)
    return config, nil
}

