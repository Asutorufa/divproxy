package MatchAndForward

import (
	"divproxy/config"
	"divproxy/net/forward"
	"divproxy/net/matcher"
	"errors"
	"fmt"
	"github.com/asticode/go-astilectron"
	"log"
	"net"
	"net/url"
)

type ForwardTo struct {
	Matcher *matcher.Match
	Config  *config.ConfigJSON
	GUI     *astilectron.Window
}

func NewForwardTo(configJsonPath, rulePath string) (forwardTo *ForwardTo, err error) {
	forwardTo = &ForwardTo{}
	forwardTo.Config, err = config.DecodeJSON(configJsonPath)
	if err != nil {
		return
	}
	forwardTo.Matcher, err = matcher.NewMatcherWithFile(forwardTo.Config.Setting.DNS, rulePath)
	return
}

func (ForwardTo *ForwardTo) log(str string) {
	if ForwardTo.GUI != nil {
		if err := ForwardTo.GUI.SendMessage(str); err != nil {
			fmt.Println(err)
		}
	} else {
		log.Println(str)
	}
}

func (ForwardTo *ForwardTo) Forward(host string) (conn net.Conn, err error) {
	var URI *url.URL
	var proxyURI *url.URL
	var proxy string
	var mode string
	if URI, err = url.Parse("//" + host); err != nil {
		return nil, err
	}
	if URI.Port() == "" {
		host = net.JoinHostPort(host, "80")
		if URI, err = url.Parse("//" + host); err != nil {
			return nil, err
		}
	}

	switch ForwardTo.Config.Setting.Bypass {
	case true:
		mode = "Bypass"
		switch ForwardTo.Matcher {
		default:
			hosts, proxy := ForwardTo.Matcher.MatchStr(URI.Hostname())
			if proxy == "block" {
				return nil, errors.New("block domain: " + host)
			}
			proxyURI = ForwardTo.Config.Nodes[proxy]
			if proxyURI == nil {
				proxyURI, err = url.Parse("direct://0.0.0.0:0")
				if err != nil {
					return nil, err
				}
			}
			for x := range hosts {
				host = net.JoinHostPort(hosts[x], URI.Port())
				conn, err = getproxyconn.ForwardTo(host, *proxyURI)
				if err == nil {
					ForwardTo.log("Mode: " + mode + " | Domain: " + host + " | match to " + proxy)
					return conn, nil
				}
			}
			return nil, errors.New("make connection:" + net.JoinHostPort(hosts[len(hosts)-1], URI.Port()) + " with proxy:" + proxy + " error")
		case nil:
			proxy = "Direct"
			URI, err = url.Parse("direct://0.0.0.0:0")
			if err != nil {
				return nil, err
			}
		}
	case false:
		switch ForwardTo.Config.Setting.Direct {
		case false:
			mode, proxy = "Only Proxy", ForwardTo.Config.Setting.Proxy
			proxyURI = ForwardTo.Config.Nodes[ForwardTo.Config.Setting.Proxy]
		case true:
			mode, proxy = "Direct", "Direct"
			proxyURI, err = url.Parse("direct://0.0.0.0:0")
			if err != nil {
				return nil, err
			}
		}
	}
	ForwardTo.log("Mode: " + mode + " | Domain: " + host + " | match to " + proxy)
	conn, err = getproxyconn.ForwardTo(host, *proxyURI)
	return
}
