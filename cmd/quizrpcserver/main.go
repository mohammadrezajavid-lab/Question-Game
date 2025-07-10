package main

import (
	"context"
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/config/quizgrpcconfig"
	"golang.project/go-fundamentals/gameapp/delivery/grpcserver/quizserver"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/infomessage"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/questionmysql"
	"golang.project/go-fundamentals/gameapp/repository/redis/redisquiz"
	"golang.project/go-fundamentals/gameapp/service/quizservice"
	"os"
	"os/signal"
)

func main() {
	var host string
	var port int
	flag.StringVar(&host, "host", "", "HTTP server host")
	flag.IntVar(&port, "port", 0, "HTTP server port")
	flag.Parse()

	grpcCfg := quizgrpcconfig.NewConfig().LoadConfig(host, port)
	logger.InitLogger(grpcCfg.LoggerCfg)

	logger.Info(fmt.Sprintf("quiz grpc config: %v", grpcCfg))

	redisAdapter := redis.New(grpcCfg.RedisCfg)
	mysqlDB := mysql.NewDB(grpcCfg.DataBaseCfg)
	setRepo := redisquiz.NewRedisDb(redisAdapter, grpcCfg.QuizRedisRepoCfg)
	dbRepo := questionmysql.NewDataBase(mysqlDB)
	quizSvc := quizservice.New(grpcCfg.QuizServiceCfg, setRepo, dbRepo)
	rpcServer := quizserver.NewQuizGrpcServer(&quizSvc, &grpcCfg.GrpcQuizCfg)

	go rpcServer.Start()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()

	logger.Info(infomessage.InfoMsgShuttingDownGracefully)

	rpcServer.Shutdown()
}
