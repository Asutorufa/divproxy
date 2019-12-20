package divproxyinit

import "divproxy/config"

func GetRuleFilePath() string {
	return "./resources/app/config/rule.config"
}

func GetConfigPath() string {
	return "./resources/app/config/config.json"
}

func GetConfig() (*config.ConfigJSON, error) {
	return config.DecodeJSON(GetConfigPath())
}
