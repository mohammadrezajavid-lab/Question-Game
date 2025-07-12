package main

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/presenceclient"
	"golang.project/go-fundamentals/gameapp/adapter/publisher"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/config/httpservercfg"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/infomessage"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/questionmysql"
	"golang.project/go-fundamentals/gameapp/repository/redis/redismatching"
	"golang.project/go-fundamentals/gameapp/repository/redis/redisset"
	"golang.project/go-fundamentals/gameapp/scheduler"
	"golang.project/go-fundamentals/gameapp/service/matchingservice"
	"golang.project/go-fundamentals/gameapp/service/quizservice"
	"os"
	"os/signal"
	"sync"
)

func main() {

	config := httpservercfg.NewConfig("", 0)

	logger.InitLogger(config.LoggerCfg)

	matchingSvc, quizSvc := setUpSvc(config)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	var wg sync.WaitGroup
	wg.Add(1)

	sch := scheduler.New(matchingSvc, quizSvc, config.SchedulerCfg)
	go sch.Start(ctx, &wg)

	<-ctx.Done()

	logger.Info(infomessage.InfoMsgShuttingDownGracefully)

	wg.Wait()

	logger.Info("Scheduler shutting down gracefully")
}

func setUpSvc(config httpservercfg.Config) (matchingservice.Service, quizservice.Service) {
	redisAdapter := redis.New(config.RedisCfg)

	presenceClient, _ := presenceclient.NewClient(config.GrpcPresenceClientCfg)
	redisPublisher := publisher.NewPublisher(config.PublisherCfg, redisAdapter)

	mysqlDB := mysql.NewDB(config.DataBaseCfg)
	setRepo := redisset.NewRedisDb(redisAdapter)
	dbRepo := questionmysql.NewDataBase(mysqlDB)

	quizSvc := quizservice.New(config.QuizServiceCfg, setRepo, dbRepo)

	matchingSvc := matchingservice.NewService(
		config.MatchingCfg,
		redismatching.NewRedisDb(redisAdapter, config.MatchingRepoCfg),
		&presenceClient,
		redisPublisher,
	)
	return matchingSvc, quizSvc
}
