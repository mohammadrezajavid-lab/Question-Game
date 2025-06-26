package main

import (
	"github.com/gorilla/websocket"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	newJwt := jwt.NewJWT(jwt.Config{
		SignKey:    "secret",
		SignMethod: "HS256",
	})

	tokenStr := newJwt.ExtractTokenFromHeader(r.Header.Get("Authorization"))
	claims, err := newJwt.ParseJWT(tokenStr)
	
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   claims.UserId,
		tokenExp: claims.ExpiresAt.UnixMicro(),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}
