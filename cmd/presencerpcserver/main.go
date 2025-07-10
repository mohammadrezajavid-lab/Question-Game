package main

import (
	"context"
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/config/presencegrpcconfig"
	"golang.project/go-fundamentals/gameapp/delivery/grpcserver/presenceserver"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/infomessage"
	"golang.project/go-fundamentals/gameapp/repository/redis/redispresence"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
	"os"
	"os/signal"
)

func main() {

	// get command
	var host string
	var port int
	flag.StringVar(&host, "host", "", "HTTP server host")
	flag.IntVar(&port, "port", 0, "HTTP server port")
	flag.Parse()

	grpcCfg := presencegrpcconfig.NewConfig().LoadConfig(host, port)

	logger.InitLogger(grpcCfg.LoggerCfg)

	logger.Info(fmt.Sprintf("grpc config: %v", grpcCfg))

	redisAdapter := redis.New(grpcCfg.RedisCfg)
	presenceSvc := presenceservice.New(redispresence.NewRedisDb(redisAdapter), grpcCfg.PresenceCfg)
	rpcServer := presenceserver.NewPresenceGrpcServer(presenceSvc, &grpcCfg.GrpcPresenceCfg)

	go rpcServer.Start()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()

	logger.Info(infomessage.InfoMsgShuttingDownGracefully)
	
	rpcServer.Shutdown()
}
