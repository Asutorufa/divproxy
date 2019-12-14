package ServerControl

import (
	"divproxy/net/MatchAndForward"
	"divproxy/net/proxy/http/server"
	"divproxy/net/proxy/socks5/server"
	"log"
)

type ServerControl struct {
	Socks5  *socks5server.ServerSocks5
	HttpS   *httpserver.HTTPServer
	forward *MatchAndForward.ForwardTo
}

func (ServerControl *ServerControl) serverControlInit() {
	var err error
	ServerControl.forward, err = MatchAndForward.NewForwardTo()
	if err != nil {
		log.Println(err)
	}
}

func (ServerControl *ServerControl) ServerStart() {
	ServerControl.serverControlInit()
	ServerControl.Socks5 = &socks5server.ServerSocks5{
		Server:    "127.0.0.1",
		Port:      "1080",
		ForwardTo: ServerControl.forward.Forward,
	}
	ServerControl.HttpS = &httpserver.HTTPServer{
		HTTPServer: "127.0.0.1",
		HTTPPort:   "8080",
		ForwardTo:  ServerControl.forward.Forward,
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
	if err = ServerControl.Socks5.Close(); err != nil {
		return
	}
	return ServerControl.HttpS.Close()
}

func (ServerControl *ServerControl) ServerRestart() {
	if err := ServerControl.ServerStop(); err != nil {
		log.Println(err)
	}
	ServerControl.ServerStart()
}
