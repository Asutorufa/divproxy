package divproxyinit

import (
	config2 "divproxy/config"
	"net/url"
)

func DeleteOneProxy(str string) error {
	configPath := GetConfigPath()
	config, err := config2.DecodeJSON(configPath)
	if err != nil {
		return err
	}
	delete(config.Nodes, str)
	if err = config2.EnCodeJSON(configPath, config); err != nil {
		return err
	}
	return nil
}

func AddOneProxy(str string, url *url.URL) error {
	configPath := GetConfigPath()
	config, err := config2.DecodeJSON(configPath)
	if err != nil {
		return err
	}
	config.Nodes[str] = url
	if err = config2.EnCodeJSON(configPath, config); err != nil {
		return err
	}
	return nil
}

func EncodeSetting(json *config2.ConfigJSON) error {
	if err := config2.EnCodeJSON(GetConfigPath(), json); err != nil {
		return err
	}
	return nil
}
