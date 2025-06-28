package websocket

import (
	"errors"
	"fmt"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
	"net/http"
)

type Config struct {
	Host           string   `mapstructure:"host"`
	Port           int      `mapstructure:"port"`
	AllowedOrigins []string `mapstructure:"allowed_origins_websocket"`
}

type WebSocket struct {
	config Config
	Hub    *Hub
	JwtCfg jwt.Config
	Server *http.Server
}

func NewWebSocket(cfg Config, jwtCfg jwt.Config) *WebSocket {
	return &WebSocket{
		config: cfg,
		Hub:    NewHub(),
		JwtCfg: jwtCfg,
		Server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		},
	}
}

func (ws *WebSocket) ServeWs() {
	router := http.NewServeMux()
	router.HandleFunc("/ws", ws.SocketHandler(ws.Hub))
	ws.Server.Handler = router

	go ws.Hub.Run()

	logger.Info(fmt.Sprintf("WebSocket Gateway started on %s", ws.Server.Addr))

	if err := ws.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err, "WebSocket Gateway server failed to start")
	}
}
