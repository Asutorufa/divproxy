package divproxyinit

import "divproxy/config"

func GetRuleFilePath() string {
	return "./rule/rule.config"
}

func GetConfigPath() string {
	return "./config/config.json"
}

func GetConfig() (*config.ConfigJSON, error) {
	return config.DecodeJSON(GetConfigPath())
}
