package MatchAndForward

import (
	"divproxy/config"
	"divproxy/net/forward"
	"divproxy/net/matcher"
	"log"
	"net"
	"net/url"
	"runtime"
)

type ForwardTo struct {
	Matcher *matcher.Match
	Config  *config.ConfigJSON
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

func (ForwardTo *ForwardTo) Forward(host string) (net.Conn, error) {
	var URI *url.URL
	var proxy string
	var mode string
	URI, err := url.Parse("//" + host)
	if err != nil {
		return nil, err
	}
	switch ForwardTo.Config.Setting.Bypass {
	case true:
		mode = "Bypass"
		switch ForwardTo.Matcher {
		default:
			hostTmp, proxy := ForwardTo.Matcher.MatchStr(URI.Hostname())
			host = hostTmp + URI.Port()
			URI = ForwardTo.Config.Nodes[proxy]
			if URI == nil {
				URI, err = url.Parse("direct://0.0.0.0:0")
				if err != nil {
					return nil, err
				}
			}
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
			URI = ForwardTo.Config.Nodes[ForwardTo.Config.Setting.Proxy]
		case true:
			mode, proxy = "Direct", "Direct"
			URI, err = url.Parse("direct://0.0.0.0:0")
			if err != nil {
				return nil, err
			}
		}
	}
	log.Println(runtime.NumGoroutine(), "Mode:", mode, "| Domain:", host, "| match to ", proxy)
	return getproxyconn.ForwardTo(host, *URI)
}
