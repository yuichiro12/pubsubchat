package main

import (
	"github.com/pebbe/zmq4"
	"log"
)

func main() {
	listener, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer listener.Close()

	err = listener.Bind("tcp://*:5001")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	publisher, err := zmq4.NewSocket(zmq4.PUB)
	err = publisher.Bind("tcp://*:5000")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	for {
		req, err := listener.Recv(0)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		log.Println(req)
		_, err = listener.Send("OK", 0)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		_, err = publisher.Send(req, 0)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}
}
