package httpserver

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// HTTPServer like name
type HTTPServer struct {
	HTTPListener *net.TCPListener
	HTTPServer   string
	HTTPPort     string
	ForwardTo    func(host string) (net.Conn, error)
	context      context.Context
	cancel       context.CancelFunc
}

func (HTTPServer *HTTPServer) httpProxyInit() error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	HTTPServer.context, HTTPServer.cancel = context.WithCancel(context.Background())
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

func (HTTPServer *HTTPServer) Close() error {
	HTTPServer.cancel()
	return HTTPServer.HTTPListener.Close()
}

func (HTTPServer *HTTPServer) httpProxyAcceptARequest() error {
	client, err := HTTPServer.HTTPListener.AcceptTCP()
	if err != nil {
		return err
	}
	if err = client.SetKeepAlivePeriod(5 * time.Second); err != nil {
		return err
	}

	go func() {
		if client == nil {
			return
		}
		defer func() {
			_ = client.Close()
		}()
		if err := HTTPServer.httpHandleClientRequest(client); err != nil {
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
	if err := HTTPServer.httpProxyInit(); err != nil {
		return err
	}
	for {
		select {
		case <-HTTPServer.context.Done():
			return nil
		default:
			if err := HTTPServer.httpProxyAcceptARequest(); err != nil {
				select {
				case <-HTTPServer.context.Done():
					return err
				default:
					log.Println(err)
					continue
				}
			}
		}
	}
}

func (HTTPServer *HTTPServer) httpHandleClientRequest2(client net.Conn) error {
	/*
		use golang http
	*/
	ss, _ := http.ReadRequest(bufio.NewReader(client))
	log.Println("header", ss.Header)
	log.Println("method", ss.Method)
	if ss.Method == http.MethodConnect {
		if _, err := client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n")); err != nil {
			return err
		}
	}
	log.Println("form", ss.Form)
	return nil
}

func (HTTPServer *HTTPServer) httpHandleClientRequest(client net.Conn) error {
	/*
		use golang http
	*/
	//ss, _ := http.ReadRequest(bufio.NewReader(client))
	//log.Println("header",ss.Header)
	//log.Println("method",ss.Method)
	//if ss.Method == http.MethodConnect{
	//	if _, err := client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n")); err != nil{
	//		return err
	//	}
	//}
	//log.Println("form",ss.Form)

	requestData := make([]byte, 1024*4)
	requestDataSize, err := client.Read(requestData[:])
	if err != nil {
		return err
	}
	if requestDataSize <= 3 {
		return nil
	}
	headerAndData := strings.Split(string(requestData[:requestDataSize]), "\r\n\r\n")
	var header, data strings.Builder
	if len(headerAndData) > 0 {
		header.WriteString(headerAndData[0])
		if len(headerAndData) > 1 {
			data.WriteString(headerAndData[1])
		}
	} else {
		return errors.New("no header")
	}

	/*
		parse request header
	*/
	headerTmp := strings.Split(header.String(), "\r\n")
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

	headerRequestSplit := strings.Split(headerTmp[0], " ")
	requestMethod := headerRequestSplit[0]
	if requestMethod == "CONNECT" {
		headerArgs["Host"] = headerRequestSplit[1]
	}
	/*
		parse request host and port
	*/
	hostPortURL, err := url.Parse("//" + headerArgs["Host"])
	if err != nil {
		return err
	}
	var headerRequest strings.Builder
	headerRequest.WriteString(strings.ReplaceAll(headerTmp[0], "http://"+hostPortURL.Host, ""))
	if hostPortURL.Port() == "" {
		hostPortURL.Host = hostPortURL.Host + ":80"
	}
	//microlog.Debug(headerArgs)
	//microlog.Debug("requestMethod:",requestMethod)
	//microlog.Debug("headerRequest ",headerRequest,"headerRequest end")
	for key, value := range headerArgs {
		headerRequest.WriteString("\r\n" + key + ": " + value)
	}
	headerRequest.WriteString("\r\n\r\n" + data.String())
	//log.Println(headerRequest.String())
	var server net.Conn
	if HTTPServer.ForwardTo != nil {
		server, err = HTTPServer.ForwardTo(hostPortURL.Host)
		if err != nil {
			return err
		}
	} else {
		server, err = net.Dial("tcp", hostPortURL.Host)
		if err != nil {
			return err
		}
	}
	defer func() {
		_ = server.Close()
	}()

	switch {
	case requestMethod == "CONNECT":
		if _, err = client.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n")); err != nil {
			return err
		}
	default:
		if _, err := server.Write([]byte(headerRequest.String())); err != nil {
			return err
		}
	}

	CloseSig := make(chan error, 0)
	go pipe(server, client, CloseSig)
	go pipe(client, server, CloseSig)
	<-CloseSig
	<-CloseSig
	close(CloseSig)
	return nil
}

func pipe(src, dst net.Conn, closeSig chan error) {
	buf := make([]byte, 0x400*32)
	for {
		n, err := src.Read(buf[0:])
		if n == 0 || err != nil {
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
