package getproxyconn

import (
	"divproxy/config"
	"divproxy/net/proxy/socks5/client"
	"net"
	"net/url"
	"time"
)

type Forward struct {
	con *config.ConfigJSON
}

func (Forward *Forward) NewForWard() (err error) {
	Forward.con, err = config.GetConfig()
	return
}

func (Forward *Forward) IsBypass(host string) (isBypass bool, proxy net.Conn) {
	if Forward.con.Setting.Bypass {
		return true, nil
	} else {
		if !Forward.con.Setting.Direct {
			proxy, _ = Forward.ForwardTo(host, *Forward.con.Nodes[Forward.con.Setting.Proxy])
		} else {
			proxy, _ = toTCP(host)
		}
	}
	return
}

func (Forward *Forward) ForwardTo(host string, proxy url.URL) (net.Conn, error) {
	if !Forward.con.Setting.Bypass {
	}
	switch proxy.Scheme {
	case "socks5":
		return toSocks5(host, proxy.Hostname(), proxy.Port())
	case "https", "http":
		return toHTTP(host, proxy.Host)
	default:
		return toTCP(host)
	}
}

func toSocks5(host string, s5Server, s5Port string) (socks5Conn net.Conn, err error) {
	return (&socks5client.Socks5Client{
		Server:           s5Server,
		Port:             s5Port,
		KeepAliveTimeout: time.Second,
		Address:          host}).NewSocks5Client()
}

func toTCP(host string) (net.Conn, error) {
	return net.Dial("tcp", host)
}

func toHTTP(host string, httpProxyServer string) (server net.Conn, err error) {
	server, err = net.Dial("tcp", httpProxyServer)
	if err != nil {
		return
	}
	if _, err = server.Write([]byte("CONNECT " + host + " HTTP/1.1\r\n\r\n")); err != nil {
		return
	}
	httpConnect := make([]byte, 1024)
	if _, err = server.Read(httpConnect[:]); err != nil {
		return
	}
	//log.Println(string(httpConnect))
	return server, nil
}
