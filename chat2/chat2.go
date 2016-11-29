package main

import (
	"log"
	"net/http"
	"github.com/pebbe/zmq4"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type Router struct {
	sessions map[string]sockjs.Session
	notifier *zmq4.Socket
}

func NewRouter() (*Router, error) {
	notifier, err := zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		return nil, err
	}
	ptr := &Router{make(map[string]sockjs.Session), notifier}
	ptr.notifier.Connect("tcp://127.0.0.1:5001")
	return ptr, nil
}

func main() {
	router, err := NewRouter()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer router.notifier.Close()

	go router.routerHandler()
	handler := sockjs.NewHandler("/chat", sockjs.DefaultOptions, router.chatHandler)
	log.Fatal(http.ListenAndServe(":8082", handler))
}

func (self *Router) chatHandler(session sockjs.Session) {
	log.Println(session.ID())
	self.sessions[session.ID()] = session
	for {
		msg, err := session.Recv()
		if err != nil {
			log.Println(session.ID(), err.Error())
			break
		}
		log.Println(">", msg)

		self.notifier.Send(msg, 0)
		self.notifier.Recv(0)
	}
	delete(self.sessions, session.ID())
}

func (self *Router) routerHandler() {
	subscriber, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		log.Fatal("n", err.Error())
	}
	defer subscriber.Close()

	err = subscriber.Connect("tcp://127.0.0.1:5000")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = subscriber.SetSubscribe("")
	if err != nil {
		log.Fatal("s", err.Error())
	}

	for {
		msg, err := subscriber.Recv(0)
		if err != nil {
			log.Fatal("r", err.Error())
		}
		log.Println("<", msg)
		// 子にメッセージを送信
		for k := range self.sessions {
			self.sessions[k].Send(msg)
		}
	}
}
