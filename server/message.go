package server

import (
	"net/http"
	"fmt"
	"text/template"
)


type Message struct {
	Url string
	index string
	Root  string
	FileHandler http.Handler
}

type Render struct {
	Key string
	Value interface{}
}



func (message *Message) ServeHTTP(res http.ResponseWriter,req *http.Request) {

	//if static file go file or some query
	if req.Method == "GET" {
		//default index
		if req.RequestURI == "" {
			req.RequestURI = "/index.html"
		}

		//show all keys
		if req.RequestURI == "/all" {
			t, err := template.ParseFiles(message.Root+"/resources/app/index.html")
			if err != nil {
				fmt.Println("parse file err:", err)
				return
			}
			for _,conf :=  range RedisHosts {
				r := GetRedis(conf.GetHval())
				r.initKeys()
				re := &Render{
					Key: "test",
					Value: r.kMap,
				}
				res.WriteHeader(200)
				if err := t.Execute(res, re); err != nil {
					res.Write([]byte(err.Error()))
					fmt.Println("There was an error:", err.Error())
				}
				return
			}
		}

		//show config router
		if req.RequestURI == "/config" {
			return
		}
		message.FileHandler.ServeHTTP(res, req)
		return
	}

	//any request must use POST
	if req.Method != "POST" {
		//do any thing
		req.ParseForm()
		DoOperation(req.PostForm)
		return
	}


	res.Write([]byte("<h1>404</h1>"))
	res.WriteHeader(404)

}

func (message *Message) getIndexContent()  {
	if message.index == "" {

	}else{

	}
}