package main

import (
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/userhandler"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/service/auth"
	"golang.project/go-fundamentals/gameapp/service/user"
	"golang.project/go-fundamentals/gameapp/validator/uservalidator"
	"log"
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

	config := httpservercfg.NewConfig(host, port)
	config.SetUpConfig(migrationCommand)
	if migrationCommand == "down" || migrationCommand == "status" {
		os.Exit(0)
	}

	//###
	fmt.Println("here")
	err := mysql.NewDB(config.DataBaseConfig).MysqlConnection.Ping()
	if err != nil {
		log.Println(err.Error())
	}

	authSvc := auth.NewService(auth.NewConfig([]byte(httpservercfg.JWTSignKey),
		httpservercfg.AccessExpirationTime,
		httpservercfg.RefreshExpirationTime,
		httpservercfg.AccessSubject,
		httpservercfg.RefreshSubject))
	userSvc := user.NewService(mysql.NewDB(config.DataBaseConfig), authSvc)
	userValidator := uservalidator.NewValidator(mysql.NewDB(config.DataBaseConfig))

	httpServer := httpserver.NewHttpServer(config, userhandler.NewHandler(userSvc, authSvc, userValidator))

	httpServer.Serve()
}
