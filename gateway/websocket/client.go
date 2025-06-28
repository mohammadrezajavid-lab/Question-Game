package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang.project/go-fundamentals/gameapp/pkg/timestamp"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   uint
	tokenExp int64
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		if timestamp.Now() > c.tokenExp {
			c.conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"token_expired"}`))
			break
		}

		// TODO - grpc call to presence upsert
		fmt.Println(msg)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for msg := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}
