package websocket

import (
	"context"
	"errors"
	"fmt"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/contract/broker"
	"golang.project/go-fundamentals/gameapp/gateway/websocket/middleware"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
	"net/http"
	"time"
)

type Config struct {
	Host                      string        `mapstructure:"host"`
	Port                      int           `mapstructure:"port"`
	AllowedOrigins            []string      `mapstructure:"allowed_origins_websocket"`
	SendBufferSize            uint          `mapstructure:"send_buffer_size"`
	BroadcastBufferSize       uint          `mapstructure:"broadcast_buffer_size"`
	GracefullyShutdownTimeout time.Duration `mapstructure:"gracefully_shutdown_timeout"`
	WebSocketPattern          string        `mapstructure:"websocket_pattern"`
	NumWorkers                uint          `mapstructure:"num_workers"`
}

type WebSocket struct {
	config         Config
	Hub            *Hub
	JwtCfg         jwt.Config
	Server         *http.Server
	presenceClient presenceclient.Client
}

func NewWebSocket(cfg Config, jwtCfg jwt.Config, presenceClientConfig presenceclient.Config, subscriber broker.Subscriber) (*WebSocket, error) {

	pClient, err := presenceclient.NewClient(presenceClientConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create presence rpc client: %w", err)
	}

	return &WebSocket{
		config: cfg,
		Hub:    NewHub(subscriber, cfg.NumWorkers, cfg.BroadcastBufferSize),
		JwtCfg: jwtCfg,
		Server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		},
		presenceClient: pClient,
	}, nil
}

func (ws *WebSocket) ServeWs() {
	router := http.NewServeMux()
	router.HandleFunc(ws.config.WebSocketPattern, ws.SocketHandler(ws.Hub))
	ws.Server.Handler = middleware.LoggerMiddleware(router)

	go ws.Hub.Run()

	logger.Info(fmt.Sprintf("WebSocket Gateway started on %s", ws.Server.Addr))

	if err := ws.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err, "WebSocket Gateway server failed to start", "addr", ws.Server.Addr)
	}
}

func (ws *WebSocket) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down WebSocket Gateway start...")

	ws.presenceClient.Close()
	logger.Info("Presence client connection closed")

	ws.Hub.Close()
	logger.Info("WebSocket Hub has been shut down")

	err := ws.Server.Shutdown(ctx)
	if err != nil {
		return err
	}

	logger.Info("WebSocket Gateway gracefully stopped")
	return nil
}
