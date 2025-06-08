package main

import (
	"context"
	"errors"
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/config/setupservices"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/infomessage"
	"golang.project/go-fundamentals/gameapp/scheduler"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

func main() {

	var host string
	var port int
	var migrationCommand string
	flag.StringVar(&host, "host", "", "HTTP server host")
	flag.IntVar(&port, "port", 0, "HTTP server port")
	flag.StringVar(
		&migrationCommand,
		"migrate-command",
		"skip",
		"Available commands are: [up] or [down] or [status] or [skip] (skip: for skipping migration for project)",
	)
	flag.Parse()

	config := httpservercfg.NewConfig(host, port)

	logger.InitLogger(config.LoggerCfg)

	config.Migrate(migrationCommand)
	if migrationCommand == "down" || migrationCommand == "status" {
		os.Exit(0)
	}

	setupSvc := setupservices.New(config)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	server := httpserver.New(
		config,
		setupSvc.AuthSvc,
		setupSvc.UserSvc,
		setupSvc.BackOfficeUserSvc,
		setupSvc.AuthorizationSvc,
		setupSvc.UserValidator,
		setupSvc.MatchingSvc,
		setupSvc.MatchingValidator,
		setupSvc.PresenceClient,
	)
	go server.Serve()

	metricsServer := &http.Server{
		Addr:    ":2112",
		Handler: promhttp.Handler(),
	}
	go func() {
		logger.Info("Starting metrics server on :2112")
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err, "Metrics server failed to start")
		}
	}()

	// start scheduler goroutine
	wg.Add(1)
	sch := scheduler.New(setupSvc.MatchingSvc, config.SchedulerCfg)
	go sch.Start(ctx, &wg)

	<-ctx.Done()

	logger.Info(infomessage.InfoMsgShuttingDownGracefully)

	// create one context, this context use for shutting down echo engine
	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.AppCfg.GracefullyShutdownTimeout)
	defer cancel()

	var shutdownWg sync.WaitGroup

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		if err := server.GetRouter().Shutdown(shutdownCtx); err != nil {
			logger.Error(err, errormessage.ErrorMsgHttpServerShutdown)
		}
	}()

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		if err := metricsServer.Shutdown(shutdownCtx); err != nil {
			logger.Error(err, "Failed to shutdown metrics server")
		}
	}()

	shutdownWg.Wait()

	wg.Wait()

	logger.Info("All services have been shut down gracefully")
}
