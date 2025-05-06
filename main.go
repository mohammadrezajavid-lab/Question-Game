package main

import (
	"flag"
	"golang.project/go-fundamentals/gameapp/config"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"os"
)

func main() {

	var host string
	var port int
	flag.StringVar(&host, "host", "127.0.0.1", "set your host for in http server")
	flag.IntVar(&port, "port", 8080, "set any port for listen http server")

	var migrationCommand string
	flag.StringVar(&migrationCommand, "migrate-command", "skip", "Available commands are: [up] or [down] or [status] or [skip] (skip: for skipping migration for project)")
	flag.Parse()

	newSetUpConfig := config.NewSetUpConfig(host, port, migrationCommand)

	if migrationCommand == "down" || migrationCommand == "status" {
		os.Exit(0)
	}

	httpServer := httpserver.NewHttpServer(newSetUpConfig.Config, newSetUpConfig.UserService, newSetUpConfig.AuthService)
	httpServer.Serve()
}
