package websocket

import (
	"golang.project/go-fundamentals/gameapp/metrics"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	quit       chan struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		quit:       make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			metrics.UserOnlineCounter.Inc()
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				metrics.UserOnlineCounter.Dec()
			}

		case <-h.quit:
			return
		}
	}
}

// Close gracefully shuts down the hub by closing all client connections
// and stopping the hub's run loop.
func (h *Hub) Close() {

	var wg sync.WaitGroup
	for client := range h.clients {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			c.Close()
		}(client)
	}

	wg.Wait()
	metrics.UserOnlineCounter.Set(0)
	close(h.quit)
}
