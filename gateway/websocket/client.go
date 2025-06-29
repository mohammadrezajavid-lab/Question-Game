package websocket

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/timestamp"
	"sync"
	"time"
)

type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	userID         uint
	tokenExp       int64
	closeOnce      sync.Once
	presenceClient presenceclient.Client
}

func (c *Client) readPump() {
	defer c.Close()

	for {
		_, msg, err := c.conn.ReadMessage()
		// TODO - fix msg
		fmt.Println(string(msg))

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Warn(err, "read data from websocket error")
			}
			break
		}

		if timestamp.Now() > c.tokenExp {
			_ = c.conn.WriteMessage(websocket.TextMessage, []byte(`{"info":"token_expired"}`))
			break
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if _, uErr := c.presenceClient.Upsert(ctx, presenceparam.NewUpsertPresenceRequest(c.userID, timestamp.Now())); uErr != nil {
			logger.Warn(uErr, fmt.Sprintf("failed to upsert presence for user user_id: %d", c.userID))
		}
	}
}

func (c *Client) writePump() {
	defer c.Close()

	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			logger.Warn(err, "write data to websocket error")
			return
		}
	}
}

// Close safely closes connection and unregisters client
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		logger.Info(fmt.Sprintf("Closing connection for user: %d", c.userID))
		c.hub.unregister <- c
		close(c.send)
		_ = c.conn.Close()
	})
}

func (c *Client) monitorTokenExpiry() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if timestamp.Now() > c.tokenExp {
				_ = c.conn.WriteMessage(websocket.TextMessage, []byte(`{"info":"token_expired"}`))
				c.Close()
				return
			}
		case <-c.hub.quit:
			return
		}
	}
}
