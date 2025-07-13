package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
	"net/http"
)

func (ws *WebSocket) NewUpgrade() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, allowed := range ws.config.AllowedOrigins {
				if origin == allowed {
					return true
				}
			}
			return false
		},
	}
}

func (ws *WebSocket) SocketHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newJwt := jwt.NewJWT(ws.JwtCfg)
		tokenStr := r.Header.Get("Authorization")

		claims, err := newJwt.ParseJWT(tokenStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}

		upgrade := ws.NewUpgrade()
		conn, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			logger.Error(err, "WebSocket upgrade failed")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "can't open websocket connection"})
			return
		}

		client := &Client{
			hub:            hub,
			conn:           conn,
			send:           make(chan []byte, ws.config.SendBufferSize),
			userID:         claims.UserId,
			tokenExp:       claims.ExpiresAt.UnixMicro(),
			presenceClient: ws.presenceClient,
		}

		hub.register <- client

		go client.writePump()
		go client.readPump()
		go client.monitorTokenExpiry()
	}
}
