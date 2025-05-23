package main

import (
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/config/setupservices"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"os"
)

func main() {

	var host string
	var port int
	flag.StringVar(&host, "host", "", "HTTP server host")
	flag.IntVar(&port, "port", 0, "HTTP server port")

	var migrationCommand string
	flag.StringVar(&migrationCommand, "migrate-command", "skip", "Available commands are: [up] or [down] or [status] or [skip] (skip: for skipping migration for project)")
	flag.Parse()

	config := httpservercfg.NewConfig(host, port)
	fmt.Println("config project: ", config)

	config.Migrate(migrationCommand)
	if migrationCommand == "down" || migrationCommand == "status" {
		os.Exit(0)
	}

	setupSvc := setupservices.New(config)
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
}
