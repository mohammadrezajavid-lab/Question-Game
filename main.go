package main

import (
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/accesscontrolmysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/usermysql"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"golang.project/go-fundamentals/gameapp/service/authorizationservice"
	"golang.project/go-fundamentals/gameapp/service/backofficeuserservice"
	"golang.project/go-fundamentals/gameapp/service/userservice"
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

	authSvc, userSvc, backOfficeUserSvc, authorizationSvc, userValidator := setupServices(config)

	httpServer := httpserver.NewHttpServer(config, authSvc, userSvc, backOfficeUserSvc, authorizationSvc, userValidator)

	httpServer.Serve()
}

func setupServices(config httpservercfg.Config) (
	*authenticationservice.Service,
	*userservice.Service,
	*backofficeuserservice.Service,
	*authorizationservice.Service,
	*uservalidator.Validator,
) {

	authSvc := authenticationservice.NewService(
		authenticationservice.NewConfig(
			config.AuthCfg.SignKey,
			config.AuthCfg.AccessExpirationTime,
			config.AuthCfg.RefreshExpirationTime,
			config.AuthCfg.AccessSubject,
			config.AuthCfg.RefreshSubject),
	)
	mysqlRepo := mysql.NewDB(config.DataBaseCfg)

	mysqlUser := usermysql.NewDataBase(mysqlRepo)
	userSvc := userservice.NewService(mysqlUser, authSvc)

	userValidator := uservalidator.NewValidator(mysqlUser)

	backOfficeUserSvc := backofficeuserservice.NewService()

	mysqlAccessControl := accesscontrolmysql.NewDataBase(mysqlRepo)
	authorizationSvc := authorizationservice.NewService(mysqlAccessControl)

	return authSvc, userSvc, backOfficeUserSvc, authorizationSvc, userValidator
}
