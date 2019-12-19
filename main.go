package main

import (
	"divproxy/gui"
)

//// Vars injected via ldflags by bundler
//var (
//	AppName            string
//	BuiltAt            string
//	VersionAstilectron string
//	VersionElectron    string
//)
//
//// Application Vars
//var (
//	debug = flag.Bool("d", true, "enables the debug mode")
//	w     *astilectron.Window
//)
func main() {
	GUI := &gui.GUI{}
	GUI.CreateGUI()

	//// Init
	//astilog.SetHandyFlags()
	//flag.Parse()
	//astilog.FlagInit()
	//
	//if err := bootstrap.Run(bootstrap.Options{
	//	Asset:    Asset,
	//	AssetDir: AssetDir,
	//	AstilectronOptions:astilectron.Options{
	//		AppName:            AppName,
	//		AppIconDarwinPath:  "resources/icon.icns",
	//		AppIconDefaultPath: "resources/icon.png",
	//		VersionAstilectron: VersionAstilectron,
	//		VersionElectron:    VersionElectron,
	//	},
	//	Debug: *debug,
	//	OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
	//		w = ws[0]
	//		go func() {
	//			time.Sleep(5 * time.Second)
	//			if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
	//				astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
	//			}
	//		}()
	//		return nil
	//	},
	//	RestoreAssets: RestoreAssets,
	//	Windows:[]*bootstrap.Window{{
	//			Homepage:"index.html",
	//		MessageHandler: nil,
	//			Options:&astilectron.WindowOptions{
	//				Center: astiptr.Bool(true),
	//				Height: astiptr.Int(600),
	//				Width:  astiptr.Int(1110),
	//			},
	//	},
	//	},
	//});err != nil{
	//	astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	//}
}
