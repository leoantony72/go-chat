package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Room struct {
	id         int
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

type Message struct {
	Message  string `json:"message,omitempty"`
	Type     string `json:"type,omitempty"`
	ClientID string `json:"client_id,omitempty"`
}

func newRoom() *Room {
	rand.Seed(time.Now().UTC().UnixNano())
	room := &Room{
		id:         rand.Int(),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	go room.run()
	return room
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			fmt.Println("client registered... room id - ", client.room.id)

			r.clients[client] = true
			fmt.Println("clients", len(r.clients))
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
			fmt.Println("clients unregistered", len(r.clients))
		case message := <-r.broadcast:
			fmt.Println(message)
			for client := range r.clients {

				client.send <- message

			}	
		}

	}
}
