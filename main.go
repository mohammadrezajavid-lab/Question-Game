package main

import (
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/config/setupservices"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"golang.project/go-fundamentals/gameapp/scheduler"
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

	go func() {
		httpServer := httpserver.NewHttpServer(
			config,
			setupSvc.AuthSvc,
			setupSvc.UserSvc,
			setupSvc.BackOfficeUserSvc,
			setupSvc.AuthorizationSvc,
			setupSvc.UserValidator,
			setupSvc.MatchingSvc,
			setupSvc.MatchingValidator,
		)
		httpServer.Serve()
	}()

	done := make(chan bool)
	go func(done <-chan bool) {
		sch := scheduler.New()
		sch.Start(done)
	}(done)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("received interrupt signal, shutting down gracefully...")
	done <- true
	time.Sleep(5 * time.Second)
}
