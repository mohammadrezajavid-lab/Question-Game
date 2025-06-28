package pprofserver

import (
	"context"
	"errors"
	"fmt"
	pprofMiddleware "golang.project/go-fundamentals/gameapp/delivery/pprofserver/middleware"
	"golang.project/go-fundamentals/gameapp/logger"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type PprofServer struct {
	config Config
	Server *http.Server
}

func NewPprofServer(cfg Config) *PprofServer {
	return &PprofServer{
		config: cfg,
		Server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:      pprofMiddleware.LogRequestMiddleware(http.DefaultServeMux),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (ps *PprofServer) Serve() {
	logger.Info(fmt.Sprintf("Starting Profiling server on %s", ps.Server.Addr))

	if err := ps.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err, "Profiling server failed to start")
	}
}

func (ps *PprofServer) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down pprof server...")

	err := ps.Server.Shutdown(ctx)
	if err != nil {
		return err
	}

	logger.Info("Pprof server gracefully stopped")
	return nil
}
