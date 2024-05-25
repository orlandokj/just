package application

import (
	"encoding/json"
)

type ApplicationConfig struct {
    Type    string `json:"type"`
    Name    string `json:"name"`
    WorkDir string `json:"workDir"`
    Config  map[string]interface{} `json:"config"`
}

func (c ApplicationConfig) GetAsString(key string) string {
    if c.Config == nil {
        return ""
    }

    value := c.Config[key]
    if value == nil {
        return ""
    }
    return value.(string)
}

func (c ApplicationConfig) ToConfigType(output interface{}) error  {
    jsonBody, err := json.Marshal(c.Config)
    if err != nil {
        return err
    }
    return json.Unmarshal(jsonBody, output)
}


