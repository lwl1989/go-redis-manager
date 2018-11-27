package main

import (
	"flag"
	"github.com/asticode/go-astilectron"
	"fmt"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astilectron-bootstrap"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

// Constants
const htmlAbout = `Welcome on <b>Astilectron</b> demo!<br>
This is using the bootstrap and the bundler.`

// Vars
var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	//w       *astilectron.Window
)

/*

 /Users/wenglong11/go/src/github.com/lwl1989/go-redis-manager/resources/vendor/electron-darwin-amd64-v1.8.1.zip
into /Users/wenglong11/go/src/github.com/lwl1989/go-redis-manager/resources/vendor/electron-darwin-amd64
failed: unzipping /Users/wenglong11/go/src/github.com/lwl1989/go-redis-manager/resources/vendor/electron-darwin-amd64-v1.8.1.zip
into /Users/wenglong11/go/src/github.com/lwl1989/go-redis-manager/resources/vendor/electron-darwin-amd64
failed: astiarchive: opening overall zip reader on /Users/wenglong11/go/src/github.com/lwl1989/go-redis-manager/resources/vendor/electron-darwin-amd64-v1.8.1.zip
failed: zip: not a valid zip file

 */
func main() {
	// Initialize astilectron
	app, _ := astilectron.New(astilectron.Options{
		AppName: "<your app name>",
		AppIconDefaultPath:   "resources/icon.png",// path is relative, it must be relative to the data directory
		AppIconDarwinPath:  "resources/icon.icns", // Same here
		BaseDirectoryPath: "resources",
	})
	defer app.Close()
	// Start astilectron
	err := app.Start()
	fmt.Println(err)
	fmt.Println("wait")
	w, _ := app.NewWindow("http://127.0.0.1:9988", &astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(600),
		Width:  astilectron.PtrInt(600),
	})
	err = w.Create()
	fmt.Println(err)
	w.OpenDevTools()
	err = w.Show()
	fmt.Println(err)
	// Blocking pattern
	app.Wait()




	return
	// Init
	flag.Parse()
	astilog.FlagInit()

	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    func(name string) ([]byte, error) {
			return []byte{}, nil
		},
		AssetDir: func(name string) ([]string, error) {
			return []string{}, nil
		},
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug: *debug,
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astilectron.PtrStr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astilectron.PtrStr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
							// Unmarshal payload
							var s string
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
								return
							}
							astilog.Infof("About modal has been displayed and payload is %s!", s)
						}); err != nil {
							astilog.Error(errors.Wrap(err, "sending about event failed"))
						}
						return
					},
				},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			go func() {
				time.Sleep(5 * time.Second)
				if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
					astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
				}
			}()
			return nil
		},
		RestoreAssets: func(dir, name string) error {
			return nil
		},
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(700),
				Width:           astilectron.PtrInt(700),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}


func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error){
	fmt.Println(m)
	return  nil,nil
}