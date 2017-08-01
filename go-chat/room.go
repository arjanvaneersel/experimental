package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"github.com/arjanvaneersel/go-chat/trace"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

type room struct {
	forward chan *message
	join chan *client
	leave chan *client
	clients map[*client]bool
	tracer trace.Tracer
}

func (rm *room) run() {
	for {
		select {
		case client := <-rm.join:
			rm.clients[client] = true
			rm.tracer.Trace("New client joined")
		case client := <-rm.leave:
			delete(rm.clients, client)
			close(client.send)
			rm.tracer.Trace("Client left")
		case msg := <-rm.forward:
			rm.tracer.Trace("Message received: ", msg.Message)
			for client := range rm.clients {
				client.send <- msg
				rm.tracer.Trace(" -- sent to client")
			}
		}
	}
}

func (rm *room) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	authCookie, err := r.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie: ", err)
		return
	}

	client := &client{
		socket: socket,
		send: make(chan *message, messageBufferSize),
		room: rm,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	rm.join <-client
	defer func() { rm.leave <- client }()
	go client.write()
	client.read()
}

func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		tracer: trace.Off(),
	}
}
