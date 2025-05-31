package main

import (
	"flag"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/config/grpcconfig"
	"golang.project/go-fundamentals/gameapp/delivery/grpcserver/presenceserver"
	"golang.project/go-fundamentals/gameapp/repository/redis/redispresence"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
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

	grpcCfg := grpcconfig.NewConfig(host, port)
	redisAdapter := redis.New(grpcCfg.RedisCfg)
	presenceSvc := presenceservice.New(redispresence.NewRedisDb(redisAdapter), grpcCfg.PresenceCfg)
	server := presenceserver.NewPresenceGrpcServer(presenceSvc)
	server.Start()
}
