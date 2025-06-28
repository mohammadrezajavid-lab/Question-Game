package main

import (
	"context"
	"flag"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/config/setupservices"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"golang.project/go-fundamentals/gameapp/delivery/metricsserver"
	"golang.project/go-fundamentals/gameapp/delivery/pprofserver"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/infomessage"
	"golang.project/go-fundamentals/gameapp/scheduler"
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

	metricServer := metricsserver.NewMetricsServer(config.MetricsCfg)
	go metricServer.Serve()

	config.Migrate(migrationCommand)
	if migrationCommand == "down" || migrationCommand == "status" {
		os.Exit(0)
	}

	setupSvc := setupservices.New(config)

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

	var profilingServer *pprofserver.PprofServer
	if config.AppCfg.DebugMod {
		profilingServer = pprofserver.NewPprofServer(config.PprofCfg)
		go profilingServer.Serve()
	}

	// start scheduler goroutine
	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	sch := scheduler.New(setupSvc.MatchingSvc, config.SchedulerCfg)
	wg.Add(1)
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
		if err := metricServer.Server.Shutdown(shutdownCtx); err != nil {
			logger.Error(err, errormessage.ErrorMsgMetricsServerShutdown)
		}
	}()

	shutdownWg.Add(1)
	go func() {
		if config.AppCfg.DebugMod {
			if err := profilingServer.Shutdown(shutdownCtx); err != nil {
				logger.Error(err, errormessage.ErrorMsgPprofServerShutdown)
			}
		}
	}()

	shutdownWg.Wait()

	wg.Wait()

	logger.Info("All services have been shut down gracefully")
}
