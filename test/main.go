package main

import (
	"divproxy/config"
	"fmt"
	"log"
	"net/url"
)

func uri(){
	URI,err := url.Parse("socks5://127.0.0.1:1080?user=xxx&password=xxx")
	if err != nil{
		log.Println(err)
	}
	fmt.Println(URI.Scheme,URI.Host,URI.Query())
}

func configTest(){
	_ = config.InitJSON("./config.json")
	s,err := config.DecodeJSON("./config.json")
	if err != nil{
		log.Println(s)
	}
	URI,err := url.Parse("socks5://127.0.0.1:1080?user=xxx&password=xxx")
	fmt.Println(URI.Port(),URI.Hostname())
	s.Nodes["xxx"] = URI
	delete(s.Nodes,"xxxx")
	_ = config.EnCodeJSON("./config.json", s)
}

func main()  {
	configTest()
}
