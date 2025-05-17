package main

import (
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/service/authentication"
	"golang.project/go-fundamentals/gameapp/service/user"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
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

	authSvc, userSvc, userValidator := setupServices(config)

	httpServer := httpserver.NewHttpServer(config, authSvc, userSvc, userValidator)

	httpServer.Serve()
}

func setupServices(config httpservercfg.Config) (*authentication.Service, *user.Service, *uservalidator.Validator) {

	authSvc := authentication.NewService(
		authentication.NewConfig(
			config.AuthCfg.SignKey,
			config.AuthCfg.AccessExpirationTime,
			config.AuthCfg.RefreshExpirationTime,
			config.AuthCfg.AccessSubject,
			config.AuthCfg.RefreshSubject),
	)

	userSvc := user.NewService(mysql.NewDB(config.DataBaseCfg), authSvc)

	userValidator := uservalidator.NewValidator(mysql.NewDB(config.DataBaseCfg))

	return authSvc, userSvc, userValidator
}
