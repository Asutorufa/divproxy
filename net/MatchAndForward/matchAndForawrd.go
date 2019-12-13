package MatchAndForward

import (
	"divproxy/config"
	"divproxy/init"
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

func NewForwardTo() (forwardTo *ForwardTo, err error) {
	forwardTo = &ForwardTo{}
	forwardTo.Config, err = divproxyinit.GetConfig()
	if err != nil {
		return
	}
	forwardTo.Matcher, err = matcher.NewMatcherWithFile(forwardTo.Config.Setting.DNS, divproxyinit.GetRuleFilePath())
	return
}

func (ForwardTo *ForwardTo) IsBypass(host string) (isBypass bool, proxy net.Conn) {
	if ForwardTo.Config.Setting.Bypass {
		return true, nil
	} else {
		if !ForwardTo.Config.Setting.Direct {
			proxy, _ = getproxyconn.ForwardTo(host, *ForwardTo.Config.Nodes[ForwardTo.Config.Setting.Proxy])
		} else {
			proxy, _ = net.Dial("tcp", host)
		}
	}
	return
}

func (ForwardTo *ForwardTo) Forward(host string) (net.Conn, error) {
	var target string
	var URI *url.URL
	var proxy string
	URI, err := url.Parse("//" + host)
	if err != nil {
		return nil, err
	}
	isBypass, server := ForwardTo.IsBypass(host)
	if isBypass {
		if ForwardTo.Matcher != nil {
			target, proxy = ForwardTo.Matcher.MatchStr(host)
			URI = ForwardTo.Config.Nodes[proxy]
			if URI == nil {
				URI, err = url.Parse("direct://0.0.0.0:0")
				if err != nil {
					return nil, err
				}
			}
		} else {
			target = host
			proxy = "direct"
			URI, err = url.Parse("direct://0.0.0.0:0")
			if err != nil {
				return nil, err
			}
		}

		server, err = getproxyconn.ForwardTo(net.JoinHostPort(target, URI.Port()), *URI)
		if err != nil {
			return nil, err
		}
	} else {
		proxy = "no bypass"
	}
	log.Println(runtime.NumGoroutine(), host, "match to", proxy)
	return server, nil
}
