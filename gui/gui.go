package gui

import (
	"divproxy/ServerControl"
	divproxyinit "divproxy/init"
	"flag"
	"fmt"
	"github.com/asticode/go-astilectron"
	_ "github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	astiptr "github.com/asticode/go-astitools/ptr"
	"github.com/pkg/errors"
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
	if gui.Window, err = gui.Accelerator.NewWindow("resources/app/index.html", &astilectron.WindowOptions{
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
	log := func(v ...interface{}) {
		if err := gui.Window.SendMessage(v); err != nil {
			fmt.Println(err)
		}
	}
	gui.server.Log = log
}

func (gui *GUI) addExtrasOptionToWindow() {
	// Open dev tools
	//_ = gui.Window.OpenDevTools()

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
		addRuleRe, err := regexp.Compile("addRule://(.*)-(.*)")
		if err != nil {
			return err
		}
		deleteRuleRe, err := regexp.Compile("deleteRule://(.*)")
		if err != nil {
			return err
		}
		applySettingRe, err := regexp.Compile("applySetting://DNS-(.*):socks5-(.*):http-(.*):proxyMode-(.*):onlyProxy-(.*)")
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
				if err := divproxyinit.AddOneProxy(name, url); err != nil {
					return err
				}
				return "add " + name + " " + url.Host + " success!"
			} else {
				return "error"
			}
		case deleteProxyRe.MatchString(s):
			s := deleteProxyRe.FindAllStringSubmatch(s, -1)
			name := s[0][1]
			fmt.Println(name)
			if err := divproxyinit.DeleteOneProxy(name); err != nil {
				return err
			}
			return "delete " + name + " success!"
		case addRuleRe.MatchString(s):
			s := addRuleRe.FindAllStringSubmatch(s, -1)
			rule := s[0][1]
			proxy := s[0][2]
			fmt.Println(rule, proxy)
			if err = divproxyinit.AddOneRule(rule+" "+proxy, divproxyinit.GetRuleFilePath()); err != nil {
				return err
			}
			return "add rule: " + rule + "-" + proxy + " success!"
		case deleteRuleRe.MatchString(s):
			s := deleteRuleRe.FindAllStringSubmatch(s, -1)
			rule := s[0][1]
			fmt.Println(rule)
			if err = divproxyinit.DeleteOneRule(rule, divproxyinit.GetRuleFilePath()); err != nil {
				return err
			}
			return "delete rule: " + rule + " success!"
		case applySettingRe.MatchString(s):
			s := applySettingRe.FindAllStringSubmatch(s, -1)
			DNS := s[0][1]
			socks5 := s[0][2]
			http := s[0][3]
			proxyMode := s[0][4]
			onlyProxy := s[0][5]
			var bypass bool
			var direct bool
			switch proxyMode {
			case "BYPASS":
				bypass = true
			case "DIRECT":
				bypass = false
				direct = true
			case "PROXY":
				bypass = false
				direct = false
			}
			fmt.Println(DNS, socks5, http, proxyMode, onlyProxy, bypass, direct)
			configs, err := divproxyinit.GetConfig()
			if err != nil {
				return err
			}
			configs.Setting.DNS = DNS
			configs.Setting.Socks5 = socks5
			configs.Setting.HTTP = http
			configs.Setting.Bypass = bypass
			configs.Setting.Direct = direct
			configs.Setting.Proxy = onlyProxy
			if err = divproxyinit.EncodeSetting(configs); err != nil {
				return err
			}
			return "apply setting!"
		}
		return nil
	})
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
