package httpserver

import (
	"errors"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"

	"divproxy/config"
	"divproxy/net/forward"
	"divproxy/net/matcher"
)

// HTTPServer like name
type HTTPServer struct {
	HTTPListener *net.TCPListener
	HTTPServer   string
	HTTPPort     string
	MatcherFile  string
	DNSServer    string
	matcher      *matcher.Match
}

func (HTTPServer *HTTPServer) HTTPProxyInit() error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	HTTPServer.matcher, err = matcher.NewMatcherWithFile(HTTPServer.DNSServer, HTTPServer.MatcherFile)
	if err != nil {
		return err
	}
	socks5ToHTTPServerIP := net.ParseIP(HTTPServer.HTTPServer)
	socks5ToHTTPServerPort, err := strconv.Atoi(HTTPServer.HTTPPort)
	if err != nil {
		return err
	}
	HTTPServer.HTTPListener, err = net.ListenTCP("tcp",
		&net.TCPAddr{IP: socks5ToHTTPServerIP, Port: socks5ToHTTPServerPort})
	if err != nil {
		return err
	}
	return nil
}

func (HTTPServer *HTTPServer) HTTPProxyAcceptARequest() error {
	HTTPConn, err := HTTPServer.HTTPListener.AcceptTCP()
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		if HTTPConn == nil {
			return
		}
		defer func() {
			_ = HTTPConn.Close()
		}()
		// log.Println("线程数:", runtime.NumGoroutine())
		err := HTTPServer.httpHandleClientRequest(HTTPConn)
		if err != nil {
			log.Println(err)
			return
		}
	}()
	return nil
}

// HTTPProxy http proxy
// server http listen server,port http listen port
// sock5Server socks5 server ip,socks5Port socks5 server port
func (HTTPServer *HTTPServer) HTTPProxy() error {
	if err := HTTPServer.HTTPProxyInit(); err != nil {
		return err
	}
	for {
		if err := HTTPServer.HTTPProxyAcceptARequest(); err != nil {
			log.Println(err)
			continue
		}
	}
}

func (HTTPServer *HTTPServer) httpHandleClientRequest(HTTPConn net.Conn) error {
	requestData := make([]byte, 1024*4)
	requestDataSize, err := HTTPConn.Read(requestData[:])
	if err != nil {
		return err
	}

	headerAndData := strings.Split(string(requestData[:requestDataSize]), "\r\n\r\n")
	var header, data string
	if len(headerAndData) > 0 {
		header = headerAndData[0]
		if len(headerAndData) > 1 {
			data = headerAndData[1]
		}
	} else {
		return errors.New("no header")
	}

	headerTmp := strings.Split(header, "\r\n")
	headerRequest := headerTmp[0]
	var requestMethod string
	headerRequestSplit := strings.Split(headerRequest, " ")
	requestMethod = headerRequestSplit[0]
	headerArgs := make(map[string]string)
	for index, line := range headerTmp {
		if index != 0 {
			//_, _ = fmt.Sscanf(line, "%s%s", &method, &host)
			tmp := strings.Split(line, ": ")
			key := tmp[0]
			value := tmp[1]
			if key == "Proxy-Connection" {
				headerArgs["Connection"] = value
				continue
			}
			headerArgs[key] = value
		}
	}

	if requestMethod == "CONNECT" {
		headerArgs["Host"] = headerRequestSplit[1]
	}
	hostPortURL, err := url.Parse("//" + headerArgs["Host"])
	if err != nil {
		return err
	}
	var address string
	if hostPortURL.Port() == "" {
		address = hostPortURL.Hostname() + ":80"
		headerRequest = strings.ReplaceAll(headerRequest, "http://"+address, "")
	} else {
		address = hostPortURL.Host
		log.Println("address:", address)
	}
	headerRequest = strings.ReplaceAll(headerRequest, "http://"+headerArgs["Host"], "")
	//microlog.Debug(headerArgs)
	//microlog.Debug("requestMethod:",requestMethod)
	//microlog.Debug("headerRequest ",headerRequest,"headerRequest end")
	//microlog.Debug("address:", address)

	for key, value := range headerArgs {
		headerRequest += "\r\n" + key + ": " + value
	}
	headerRequest += "\r\n\r\n" + data

	var domainPort string
	if net.ParseIP(hostPortURL.Hostname()) == nil {
		domainPort = strings.Split(address, ":")[1]
	} else if net.ParseIP(hostPortURL.Hostname()).To4() != nil {
		domainPort = strings.Split(address, ":")[1]
	} else {
		domainPort = strings.Split(address, "]:")[1]
	}

	log.Println(hostPortURL.Hostname())
	_, proxy := HTTPServer.matcher.MatchStr(hostPortURL.Hostname())
	s, err := config.GetConfig()
	if err != nil {
		return err
	}
	Conn, err := getproxyconn.Forward(net.JoinHostPort(hostPortURL.Hostname(), domainPort), *s.Nodes[proxy])
	if err != nil{
		return err
	}
	defer func() {
		if err = Conn.Close(); err !=nil{
			log.Println(err)
		}
	}()

	switch {
	case requestMethod == "CONNECT":
		if _, err = HTTPConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n")); err != nil{
			return err
		}
	default:
		if _, err := Conn.Write([]byte(headerRequest)); err != nil {
			return err
		}
	}

	ConnToHTTPConnCloseSig, HTTPConnToConnCloseSig := make(chan error, 1), make(chan error, 1)
	go pipe(Conn, HTTPConn, ConnToHTTPConnCloseSig)
	go pipe(HTTPConn, Conn, HTTPConnToConnCloseSig)
	<-ConnToHTTPConnCloseSig
	close(ConnToHTTPConnCloseSig)
	<-HTTPConnToConnCloseSig
	close(HTTPConnToConnCloseSig)
	return nil
}

func pipe(src, dst net.Conn, closeSig chan error) {
	buf := make([]byte, 0x400*32)
	for {
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
