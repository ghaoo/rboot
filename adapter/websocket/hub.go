package websocket

import (
	"encoding/json"
	"github.com/ghaoo/rboot"
	"net/http"
)

// Hub 维护活动的客户端并向客户端广播消息
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client

	incoming chan *rboot.Message
	outgoing chan *rboot.Message
}

func newHub(bot *rboot.Robot) rboot.Adapter {
	hub := &Hub{
		incoming:   make(chan *rboot.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	go hub.run()

	bot.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	return hub
}

func (h *Hub) Incoming() chan *rboot.Message {
	return h.incoming
}

func (h *Hub) Outgoing() chan *rboot.Message {
	return h.outgoing
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.outgoing:
			msg, _ := json.Marshal(message)
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
