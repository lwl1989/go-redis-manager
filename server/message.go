package server

import (
	"net/http"
	"fmt"
)


type Message struct {
	Url string
	index string
	Root  string
	FileHandler http.Handler
}


func (message *Message) ServeHTTP(res http.ResponseWriter,req *http.Request) {

	//if static file go file
	if req.Method == "GET" {
		fmt.Println(req)
		if req.RequestURI == "" {
			req.RequestURI = "/index.html"
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