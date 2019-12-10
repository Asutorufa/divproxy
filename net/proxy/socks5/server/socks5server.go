package socks5server

import (
	"context"
	"divproxy/config"
	"errors"
	"log"
	"net"
	"runtime"
	"strconv"

	"divproxy/net/forward"
	"divproxy/net/matcher"
)

// ServerSocks5 <--
type ServerSocks5 struct {
	Server           string
	Port             string
	Username         string
	Password         string
	MatcherFile      string
	DNSServer        string
	Context          context.Context
	ContextCloseConn context.Context
	Cancel           context.CancelFunc
	CancelCloseConn  context.CancelFunc
	conn             *net.TCPListener
	matcher          *matcher.Match
}

func (socks5Server *ServerSocks5) socks5Init() error {
	socks5Server.Context, socks5Server.Cancel = context.WithCancel(context.Background())
	socks5Server.ContextCloseConn, socks5Server.CancelCloseConn = context.WithCancel(context.Background())
	var err error
	if socks5Server.MatcherFile != "" {
		socks5Server.matcher, err = matcher.NewMatcherWithFile(socks5Server.DNSServer, socks5Server.MatcherFile)
		if err != nil {
			return err
		}
	} else {
		socks5Server.matcher = matcher.NewMatcher(socks5Server.DNSServer)
	}
	socks5ServerIP := net.ParseIP(socks5Server.Server)
	socks5ServerPort, err := strconv.Atoi(socks5Server.Port)
	if err != nil {
		return err
	}
	socks5Server.conn, err = net.ListenTCP("tcp", &net.TCPAddr{IP: socks5ServerIP, Port: socks5ServerPort})
	if err != nil {
		return err
	}
	return nil
}

func (socks5Server *ServerSocks5) socks5AcceptARequest() error {
	client, err := socks5Server.conn.AcceptTCP()
	if err != nil {
		return err
	}

	go func() {
		if client == nil {
			return
		}
		defer func() {
			_ = client.Close()
		}()
		socks5Server.handleClientRequest(client)
	}()
	return nil
}

// Socks5 <--
func (socks5Server *ServerSocks5) Socks5() error {
	if err := socks5Server.socks5Init(); err != nil {
		return err
	}
	for {
		select {
		case <-socks5Server.Context.Done():
			if err := socks5Server.conn.Close(); err != nil {
				log.Println(err)
			}
			socks5Server.CancelCloseConn()
			return errors.New("socks5 server close")
		default:
			if err := socks5Server.socks5AcceptARequest(); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (socks5Server *ServerSocks5) handleClientRequest(client net.Conn) {
	var b [1024]byte
	_, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}

	if b[0] == 0x05 { //只处理Socks5协议
		_, _ = client.Write([]byte{0x05, 0x00})
		if b[1] == 0x01 {
			// 对用户名密码进行判断
			if b[2] == 0x02 {
				_, err = client.Read(b[:])
				if err != nil {
					log.Println(err)
					return
				}
				username := b[2 : 2+b[1]]
				password := b[3+b[1] : 3+b[1]+b[2+b[1]]]
				if socks5Server.Username == string(username) && socks5Server.Password == string(password) {
					_, _ = client.Write([]byte{0x01, 0x00})
				} else {
					_, _ = client.Write([]byte{0x01, 0x01})
					return
				}
			}
		}

		n, err := client.Read(b[:])
		if err != nil {
			log.Println(err)
			return
		}

		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		switch b[1] {
		case 0x01:
			target, proxy := socks5Server.matcher.MatchStr(host)
			s, err := config.GetConfig()
			if err != nil {
				log.Println(err)
				return
			}
			url := s.Nodes[proxy]
			if url == nil {
				url, err = url.Parse("direct://0.0.0.0:0")
			}
			server, err := getproxyconn.Forward(net.JoinHostPort(target, port), *url)
			if err != nil {
				log.Println(err)
				return
			}
			_, _ = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //respond to connect successful
			forward(client, server, socks5Server.Context)

		case 0x02:
			log.Println("bind 请求 " + net.JoinHostPort(host, port))

		case 0x03:
			log.Println("udp 请求 " + net.JoinHostPort(host, port))
			socks5Server.udp(client, net.JoinHostPort(host, port))
		}
	}
}

func (socks5Server *ServerSocks5) connect() {
	// do something
}

func (socks5Server *ServerSocks5) udp(client net.Conn, domain string) {
	server, err := net.Dial("udp", domain)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = server.Close()
	}()
	_, _ = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //respond to connect successful

	// forward
	forward(server, client, socks5Server.Context)
}

func forward(src, dst net.Conn, ctx context.Context) {
	CloseSig := make(chan error, 0)
	go pipe(src, dst, CloseSig, ctx)
	go pipe(dst, src, CloseSig, ctx)
	<-CloseSig
	log.Println(runtime.NumGoroutine(), "close")
}

func pipe(src, dst net.Conn, closeSig chan error, ctx context.Context) {
	buf := make([]byte, 0x400*32)
	for {
		select {
		case <-ctx.Done():
			log.Println("close forward")
			_ = src.Close()
			_ = dst.Close()
			closeSig <- nil
			return
		default:
			n, err := src.Read(buf[0:])
			if err != nil {
				closeSig <- err
				return
			}
			_, err = dst.Write(buf[0:n])
			if err != nil {
				closeSig <- err
				return
			}
		}
	}
}
