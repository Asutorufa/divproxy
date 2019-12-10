package config

import (
	"encoding/json"
	"net/url"
	"os"
)

// configJSON config json struct
type configJSON struct {
	Nodes   map[string]*url.URL `json:"nodes"`
}

// InitJSON init the config json file
func InitJSON(configPath string) (err error) {
	return EnCodeJSON(configPath, &configJSON{
		Nodes:   map[string]*url.URL{},
	})
}

func DecodeJSON(configPath string) (pa *configJSON,err error) {
	pa = &configJSON{}
	file, err := os.Open(configPath)
	if err != nil {
		return
	}
	if err = json.NewDecoder(file).Decode(&pa); err != nil {
		return
	}
	return
}

func EnCodeJSON(configPath string, pa *configJSON) (err error) {
	file, err := os.Create(configPath)
	if err != nil {
		return
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	if err := enc.Encode(&pa); err != nil {
		return
	}
	return
}

func GetConfig()  (*configJSON, error) {
	return DecodeJSON("./config/config.json")
}