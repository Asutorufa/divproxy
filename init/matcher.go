package divproxyinit

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
	forward *getproxyconn.Forward
}

func GetMatcher() (match *matcher.Match) {
	match, err := matcher.NewMatcherWithFile("114.114.114.114:53", "./rule/rule.config")
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func NewForwardTo() (forwardTo *ForwardTo, err error) {
	forwardTo = &ForwardTo{}
	forwardTo.Matcher = GetMatcher()
	if forwardTo.forward, err = getproxyconn.NewForWard(); err != nil {
		return
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
	isBypass, server := ForwardTo.forward.IsBypass(host)
	if isBypass {
		if ForwardTo.Matcher != nil {
			target, proxy = ForwardTo.Matcher.MatchStr(host)
			s, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			URI = s.Nodes[proxy]
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

		server, err = ForwardTo.forward.ForwardTo(net.JoinHostPort(target, URI.Port()), *URI)
		if err != nil {
			return nil, err
		}
	} else {
		proxy = "no bypass"
	}
	log.Println(runtime.NumGoroutine(), host, "match to", proxy)
	return server, nil
}
