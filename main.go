package main

import (
	"golang.project/go-fundamentals/gameapp/config"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
)

func main() {

	newSetUpConfig := config.NewSetUpConfig()

	httpServer := httpserver.NewHttpServer(newSetUpConfig.Config, newSetUpConfig.UserService, newSetUpConfig.AuthService)
	httpServer.Serve()
}
