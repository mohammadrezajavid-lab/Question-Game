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

	// TODO - It's better to separate setUpService and setUpValidator from setUpConfig
	setUpConfig := config.NewSetUpConfig(host, port, migrationCommand)
	//setUpService :=
	//setUpValidator :=

	if migrationCommand == "down" || migrationCommand == "status" {
		os.Exit(0)
	}

	httpServer := httpserver.NewHttpServer(
		setUpConfig.Config,
		setUpConfig.UserService,
		setUpConfig.AuthService,
		setUpConfig.UserValidator,
	)

	httpServer.Serve()
}
