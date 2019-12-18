package gui

import (
	"divproxy/ServerControl"
	divproxyinit "divproxy/init"
	"flag"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	astiptr "github.com/asticode/go-astitools/ptr"
	"github.com/pkg/errors"
	"io/ioutil"
	url2 "net/url"
	"regexp"
)

type GUI struct {
	Accelerator *astilectron.Astilectron
	Window      *astilectron.Window
	server      *ServerControl.ServerControl
}

func (gui *GUI) guiInit() {
	// Parse flags
	astilog.SetHandyFlags()
	flag.Parse()
	astilog.FlagInit()

	// Create astilectron
	var err error
	gui.Accelerator, err = astilectron.New(astilectron.Options{AppName: "Test"})
	if err != nil {
		astilog.Fatal(errors.Wrap(err, "main: creating astilectron failed"))
	}
	//defer gui.Accelerator.Close()

	// Handle signals
	gui.Accelerator.HandleSignals()

	// Start
	if err = gui.Accelerator.Start(); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: starting astilectron failed"))
	}
	gui.server = &ServerControl.ServerControl{
		ConfigJsonPath: divproxyinit.GetConfigPath(),
		RulePath:       divproxyinit.GetRuleFilePath(),
	}
}

func (gui *GUI) createWindow() {
	// New window
	var err error
	if gui.Window, err = gui.Accelerator.NewWindow("src/index.html", &astilectron.WindowOptions{
		Center: astiptr.Bool(true),
		Height: astiptr.Int(600),
		Width:  astiptr.Int(1110),
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: new window failed"))
	}

	// Create windows
	if err = gui.Window.Create(); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: creating window failed"))
	}
	gui.server.GUI = gui.Window
}

func (gui *GUI) addExtrasOptionToWindow() {
	// Open dev tools
	_ = gui.Window.OpenDevTools()

	// Close dev tools
	//w.CloseDevTools()
}

func (gui *GUI) addListenerToWindows() {
	// This will send a message and execute a callback
	// Callbacks are optional
	//_ = gui.Window.SendMessage("hello", func(m *astilectron.EventMessage) {
	//	// Unmarshal
	//	var s string
	//	_ = m.Unmarshal(&s)
	//
	//	// Process message
	//	astilog.Debugf("received %s", s)
	//	fmt.Println(s)
	//})
	//_ = w.SendMessage("hello")

	// This will listen to messages sent by Javascript
	gui.Window.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		_ = m.Unmarshal(&s)

		addProxyRe, err := regexp.Compile("addProxy://(.*)-(.*)")
		if err != nil {
			return err
		}
		deleteProxyRe, err := regexp.Compile("deleteProxy://(.*)")
		if err != nil {
			return err
		}
		getRuleRe, err := regexp.Compile("getRule://")
		if err != nil {
			return err
		}

		switch {
		case s == "startProxy://":
			gui.server.ServerStart()
			return "start"
		case s == "restartProxy://":
			gui.server.ServerRestart()
			return "restart"
		case s == "stopProxy://":
			if err := gui.server.ServerStop(); err != nil {
				return err
			}
			return "stop"
		case addProxyRe.MatchString(s):
			s := addProxyRe.FindAllStringSubmatch(s, -1)
			if len(s[0]) >= 3 {
				name := s[0][1]
				url, err := url2.Parse(s[0][2])
				if err != nil {
					return err
				}
				fmt.Println(name, url)
				return "add " + name + " " + url.Host + " success!"
			} else {
				return "error"
			}
		case deleteProxyRe.MatchString(s):
			s := deleteProxyRe.FindAllStringSubmatch(s, -1)
			name := s[0][1]
			fmt.Println(name)
			return "delete " + name + " success!"
		case getRuleRe.MatchString(s):
			configTemp, err := ioutil.ReadFile("./rule/rule.config")
			if err != nil {
				return err
			}
			return string(configTemp)
		}
		return nil
	})
}

func (gui *GUI) Log(str string) {
	if err := gui.Window.SendMessage("hello"); err != nil {
		fmt.Println(err)
	}
}

func (gui *GUI) CreateGUI() {
	gui.guiInit()
	gui.createWindow()
	gui.addExtrasOptionToWindow()
	gui.addListenerToWindows()
	// Blocking pattern
	defer gui.Accelerator.Close()
	gui.Accelerator.Wait()
}
