package main

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	http.HandleFunc("/", handler)
	http.Serve(listen, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Fprintf(w, string(file[:]))
}
