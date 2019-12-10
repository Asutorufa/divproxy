package main

import (
	"flag"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/ptr"
	"github.com/pkg/errors"
	"io/ioutil"
	"regexp"
)

func main() {
	// Parse flags
	astilog.SetHandyFlags()
	flag.Parse()
	astilog.FlagInit()

	// Create astilectron
	a, err := astilectron.New(astilectron.Options{AppName: "Test"})
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "main: creating astilectron failed"))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: starting astilectron failed"))
	}

	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow("src/index.html", &astilectron.WindowOptions{
		Center:          astiptr.Bool(true),
		Height:          astiptr.Int(600),
		Width:           astiptr.Int(1110),
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: new window failed"))
	}

	// Create windows
	if err = w.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: creating window failed"))
	}

	// Open dev tools
	_ = w.OpenDevTools()

	// Close dev tools
	//w.CloseDevTools()

	// This will send a message and execute a callback
	// Callbacks are optional
	_ = w.SendMessage("hello", func(m *astilectron.EventMessage) {
		// Unmarshal
		var s string
		_ = m.Unmarshal(&s)

		// Process message
		astilog.Debugf("received %s", s)
		fmt.Println(s)
	})
	//_ = w.SendMessage("hello")

	// This will listen to messages sent by Javascript
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		_ = m.Unmarshal(&s)

		addProxyRe, err := regexp.Compile("addProxy://(.*)-(.*)")
		if err != nil{
			return err
		}
		deleteProxyRe, err := regexp.Compile("deleteProxy://(.*)")
		if err != nil{
			return err
		}

		getRuleRe, err := regexp.Compile("getRule://")
		if err != nil{
			return err
		}

		switch {
		case addProxyRe.MatchString(s):
			s := addProxyRe.FindAllStringSubmatch(s, -1)
			fmt.Println(s)
			return "add success!"
		case deleteProxyRe.MatchString(s):
			s := deleteProxyRe.FindAllStringSubmatch(s,-1)
			fmt.Println(s)
			return "delete success!"
		case getRuleRe.MatchString(s):
			configTemp, err := ioutil.ReadFile("./rule/rule.config")
			if err != nil{
				return err
			}
			return string(configTemp)
		case s == "hello":
			//_ = w.SendMessage("hello")
			configTemp, _ := ioutil.ReadFile("./config/config.json")
			//fmt.Println(configTemp)
			return string(configTemp)
		}
		return nil
	})

	// Blocking pattern
	a.Wait()
}