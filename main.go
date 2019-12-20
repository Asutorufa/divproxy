package main

import "divproxy/gui"

import divproxyinit "divproxy/init"

import "log"

func main() {
    if err := divproxyinit.InitConfig(); err != nil{
        log.Println(err)
        return
    }
	GUI := &gui.GUI{}
	GUI.CreateGUI()
}
