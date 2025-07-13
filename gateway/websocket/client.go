package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/timestamp"
	"strconv"
	"strings"
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

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				metrics.FailedReadPompCounter.Inc()
				logger.Error(err, "Read data from websocket error", "user_id", c.userID)
			}
			break
		}

		if timestamp.Now() > c.tokenExp {
			_ = c.conn.WriteMessage(websocket.TextMessage, []byte(`{"info":"token_expired"}`))
			break
		}

		if strings.TrimSpace(string(msg)) != `{"event":"heartbeat"}` {
			metrics.InvalidHeartBeatMessageCounter.With(prometheus.Labels{"user_id": strconv.Itoa(int(c.userID))}).Inc()
			logger.Info("Invalid heartbeat message", "user_id", c.userID, "message", string(msg))
			continue
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if _, uErr := c.presenceClient.Upsert(ctx, presenceparam.NewUpsertPresenceRequest(c.userID, timestamp.Now())); uErr != nil {
			logger.Error(uErr, "Failed Upsert presence for user", "user_id", c.userID)
		}

		logger.Info("received heartbeat message", "user_id", c.userID)
	}
}

func (c *Client) writePump() {
	defer c.Close()

	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			metrics.FailedWritePompCounter.Inc()
			logger.Error(err, "write data to websocket error", "user_id", c.userID)
			return
		}
	}
}

// Close safely closes connection and unregisters client
func (c *Client) Close() {
	c.closeOnce.Do(
		func() {
			metrics.ClosingConnectionWebsocket.Inc()
			logger.Info("Closing connection websocket.", "user_id", c.userID)
			c.hub.unregister <- c
			_ = c.conn.Close()
		},
	)
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
