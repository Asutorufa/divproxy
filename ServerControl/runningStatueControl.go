package ServerControl

import (
	"divproxy/MatchAndForward"
	config2 "divproxy/config"
	httpserver "divproxy/net/proxy/http/server"
	socks5server "divproxy/net/proxy/socks5/server"
	"errors"
	"fmt"
	"log"
	"net/url"
)

type ServerControl struct {
	Socks5         *socks5server.ServerSocks5
	HttpS          *httpserver.HTTPServer
	forward        *MatchAndForward.ForwardTo
	config         *config2.ConfigJSON
	Log            func(v ...interface{})
	ConfigJsonPath string
	RulePath       string
}

func (ServerControl *ServerControl) serverControlInit() {
	var err error
	ServerControl.forward, err = MatchAndForward.NewForwardTo(ServerControl.ConfigJsonPath, ServerControl.RulePath)
	if err != nil {
		log.Println(err)
	}
	ServerControl.forward.Log = ServerControl.Log
	if ServerControl.config, err = config2.DecodeJSON(ServerControl.ConfigJsonPath); err != nil {
		log.Println(err)
	}
}

func (ServerControl *ServerControl) ServerStart() {
	ServerControl.serverControlInit()
	var err error
	socks5, err := url.Parse("//" + ServerControl.config.Setting.Socks5)
	if err != nil {
		log.Println(err)
	}
	ServerControl.Socks5, err = socks5server.NewSocks5Server(socks5.Hostname(), socks5.Port(), "", "", ServerControl.forward.Forward)
	if err != nil {
		log.Println(err)
	}
	http, err := url.Parse("//" + ServerControl.config.Setting.HTTP)
	if err != nil {
		log.Println(err)
	}
	ServerControl.HttpS, err = httpserver.NewHTTPServer(http.Hostname(), http.Port(), "", "", ServerControl.forward.Forward)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		if err := ServerControl.Socks5.Socks5(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		if err := ServerControl.HttpS.HTTPProxy(); err != nil {
			log.Println(err)
		}
	}()
}

func (ServerControl *ServerControl) ServerStop() (err error) {
	if ServerControl.Socks5 != nil {
		if err = ServerControl.Socks5.Close(); err != nil {
			return
		}
	}
	if ServerControl.HttpS != nil {
		if err = ServerControl.HttpS.Close(); err != nil {
			return
		}
	}
	if ServerControl.Socks5 != nil && ServerControl.HttpS != nil {
		ServerControl.HttpS = nil
		ServerControl.Socks5 = nil
		return nil
	}
	return errors.New("not Start")
}

func (ServerControl *ServerControl) ServerRestart() {
	if err := ServerControl.ServerStop(); err != nil {
		fmt.Println(err)
	} else {
		ServerControl.ServerStart()
	}
}
