package config

import (
	"encoding/json"
	"net/url"
	"os"
)

type settingJSON struct {
	DNS    string `json:"dns"`
	Socks5 string `json:"socks_5"`
	HTTP   string `json:"http"`
	Bypass bool   `json:"bypass"`
	Direct bool   `json:"direct"`
	Proxy  string `json:"proxy"`
}

// configJSON config json struct
type ConfigJSON struct {
	Nodes   map[string]*url.URL `json:"nodes"`
	Setting *settingJSON        `json:"setting"`
}

// InitJSON init the config json file
func InitJSON(configPath string) (err error) {
	return EnCodeJSON(configPath, &ConfigJSON{
		Nodes:   map[string]*url.URL{},
		Setting: &settingJSON{},
	})
}

func DecodeJSON(configPath string) (pa *ConfigJSON, err error) {
	pa = &ConfigJSON{}
	file, err := os.Open(configPath)
	if err != nil {
		return
	}
	if err = json.NewDecoder(file).Decode(&pa); err != nil {
		return
	}
	return
}

func EnCodeJSON(configPath string, pa *ConfigJSON) (err error) {
	file, err := os.Create(configPath)
	if err != nil {
		return
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	if err = enc.Encode(&pa); err != nil {
		return
	}
	return
}

func GetConfig() (*ConfigJSON, error) {
	return DecodeJSON("./config/config.json")
}
