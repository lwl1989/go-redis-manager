package main

import (
	"github.com/asticode/go-astilectron"
	"fmt"
	"github.com/asticode/go-astilog"
	"time"
	"net/url"
	"net/http"
	"path/filepath"
	"os"
	"github.com/lwl1989/go-redis-manager/server"
)

// Constants
const htmlAbout = `Welcome on <b>Astilectron</b> demo!<br>
This is using the bootstrap and the bundler.`

var debug bool

func main() {
	// Initialize astilectron
	debug = true
	urlStr := "http://127.0.0.1:10003"
	conf := make(map[string]*server.RedisConfig)
	redisConfg := &server.RedisConfig{
		Host:"127.0.0.1:6379",
		Pw:"",
		Db:0,
	}
	conf[redisConfg.GetHval()] = redisConfg
	server.RedisHosts = conf
	buildHttpServerHandler(urlStr)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//buildEctron(urlStr)
	//select {
	//
	//}
	return
}

func buildHttpServerHandler(urlStr string) (err error) {

	urlObj,err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	var root string
	if !debug {
		root, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}
	}else{
		root = os.Getenv("GOPATH")+"/src/github.com/lwl1989/go-redis-manager"
	}
	//fmt.Println(root)
	message := &server.Message{
		Url: urlStr,
		Root: root,
		FileHandler:http.FileServer(http.Dir(root+"/resources/app")),
	}
	//fmt.Println(urlObj.Port())
	//+urlObj.Port()
	err = http.ListenAndServe(":"+urlObj.Port(), message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func buildEctron(urlStr string) {
	app, err := astilectron.New(astilectron.Options{
		AppName: "oduoke",
		AppIconDefaultPath:   "icon.png",// path is relative, it must be relative to the data directory
		AppIconDarwinPath:  "icon.icns", // Same here
		BaseDirectoryPath: "resources",
		AcceptTCPTimeout: 3*time.Second,
	})
	defer app.Close()
	fmt.Println(err)
	// Start astilectron
	err = app.Start()
	//fmt.Println(err)
	//fmt.Println(app)
	w, err := app.NewWindow(urlStr, &astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(600),
		Width:  astilectron.PtrInt(600),
	})


	app.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
		astilog.Error("App has crashed")
		return
	})

	// Add a listener on the window
	w.On(astilectron.EventNameWindowEventResize, func(e astilectron.Event) (deleteListener bool) {
		astilog.Info("Window resized")
		return
	})

	fmt.Println(err)
	err = w.Create()
	w.OpenDevTools()
	err = w.Show()
	// Blocking pattern
	app.Wait()
}