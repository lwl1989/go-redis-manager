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

	//if static file go file
	if req.Method == "GET" {
		fmt.Println(req)
		if req.RequestURI == "" {
			req.RequestURI = "/index.html"
		}
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
					Value: r.kMap.String(),
				}
				res.WriteHeader(200)
				if err := t.Execute(res, re); err != nil {
					res.Write([]byte(err.Error()))
					fmt.Println("There was an error:", err.Error())
				}


				return
			}

		}
		message.FileHandler.ServeHTTP(res, req)
		return
	}

	//any request must use POST
	if req.Method != "POST" {
		res.Write([]byte("<h1>500</h1>"))
		res.WriteHeader(500)
		return
	}

	//do any thing
}

func (message *Message) getIndexContent()  {
	if message.index == "" {

	}else{

	}
}