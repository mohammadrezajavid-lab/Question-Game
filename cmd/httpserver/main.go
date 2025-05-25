package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/config/setupservices"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"os"
	"os/signal"
	"time"
)

func main() {

	var host string
	var port int
	flag.StringVar(&host, "host", "", "HTTP server host")
	flag.IntVar(&port, "port", 0, "HTTP server port")

	var migrationCommand string
	flag.StringVar(
		&migrationCommand,
		"migrate-command",
		"skip",
		"Available commands are: [up] or [down] or [status] or [skip] (skip: for skipping migration for project)",
	)
	flag.Parse()

	config := httpservercfg.NewConfig(host, port)
	fmt.Println("config project: ", config)

	config.Migrate(migrationCommand)
	if migrationCommand == "down" || migrationCommand == "status" {
		os.Exit(0)
	}

	setupSvc := setupservices.New(config)

	var httpServerChan *echo.Echo

	go func() {
		httpServer := httpserver.New(
			config,
			setupSvc.AuthSvc,
			setupSvc.UserSvc,
			setupSvc.BackOfficeUserSvc,
			setupSvc.AuthorizationSvc,
			setupSvc.UserValidator,
			setupSvc.MatchingSvc,
			setupSvc.MatchingValidator,
		)

		httpServerChan = httpServer.Serve()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("received interrupt signal, shutting down gracefully...")

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, config.AppCfg.GracefullyShutdownTimeout)
	defer cancel()

	if sErr := httpServerChan.Shutdown(ctxWithTimeout); sErr != nil {
		fmt.Printf("\nhttp server shutdown error: %v\n", sErr)
	}

	time.Sleep(config.AppCfg.GracefullyShutdownTimeout)
}
