package config

import (
	"net/url"
	"testing"
)

func TestConfig(t *testing.T) {
	if err := InitJSON("config.json"); err != nil {
		t.Error(err)
		return
	}
	url1, err := url.Parse("socks5://127.0.0.1:1080?user=xxx&password=xxx")
	if err != nil {
		t.Error(err)
		return
	}
	url2, err := url.Parse("HTTP://127.0.0.1:8080")
	if err != nil {
		t.Error(err)
		return
	}
	config, err := DecodeJSON("config.json")
	if err != nil {
		t.Error(err)
		return
	}
	config.Nodes["url1"] = url1
	config.Nodes["url2"] = url2
	config.Setting.DNS = "8.8.8.8:53"
	config.Setting.Socks5 = "127.0.0.1:1081"
	config.Setting.HTTP = "127.0.0.1:8081"
	config.Setting.Bypass = true
	config.Setting.Direct = false
	if err = EnCodeJSON("config.json", config); err != nil {
		t.Error(err)
		return
	}
}
