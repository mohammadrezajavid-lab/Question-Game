package main

import (
	"context"
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/config/setupservices"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"golang.project/go-fundamentals/gameapp/scheduler"
	"os"
	"os/signal"
	"sync"
)

func main() {

	// get command
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

	// setup http server config
	config := httpservercfg.NewConfig(host, port)
	fmt.Println("config project: ", config)

	// run migrations
	config.Migrate(migrationCommand)
	if migrationCommand == "down" || migrationCommand == "status" {
		os.Exit(0)
	}

	// setup services
	setupSvc := setupservices.New(config)

	// start http server goroutine
	server := httpserver.New(
		config,
		setupSvc.AuthSvc,
		setupSvc.UserSvc,
		setupSvc.BackOfficeUserSvc,
		setupSvc.AuthorizationSvc,
		setupSvc.UserValidator,
		setupSvc.MatchingSvc,
		setupSvc.MatchingValidator,
	)
	go server.Serve()

	// start scheduler goroutine
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	sch := scheduler.New()
	go sch.Start(done, &wg)

	// waiting for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit // blocked main in this line
	fmt.Println("received interrupt signal, shutting down gracefully...")

	// create one context, this context use for shutting down echo engine
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, config.AppCfg.GracefullyShutdownTimeout)
	defer cancel()

	// shutdown echo engine
	if sErr := server.GetRouter().Shutdown(ctxWithTimeout); sErr != nil {
		fmt.Printf("\nhttp server shutdown error: %v\n", sErr)
	}

	// stopping scheduler
	done <- true

	// waiting for stop scheduler
	wg.Wait()

	<-ctxWithTimeout.Done()
}
