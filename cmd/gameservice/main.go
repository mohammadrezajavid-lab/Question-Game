package main

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/adapter/publisher"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/adapter/subscriber"
	"golang.project/go-fundamentals/gameapp/config/gameservicecfg"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/infomessage"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/repository/mysql/gamemysql"
	"golang.project/go-fundamentals/gameapp/service/gameservice"
	"os"
	"os/signal"
	"sync"
)

func main() {

	config := gameservicecfg.NewConfig().LoadConfig()
	fmt.Println(config)

	logger.InitLogger(config.LoggerCfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	game := setUpSvc(config)

	wg.Add(1)
	go game.Start(ctx, &wg)

	<-ctx.Done()

	logger.Info(infomessage.InfoMsgShuttingDownGracefully)

	wg.Wait()

	logger.Info("Game Service shutting down gracefully")
}

func setUpSvc(config *gameservicecfg.Config) gameservice.Service {
	redisAdapter := redis.New(config.RedisCfg)
	mysqlDB := mysql.NewDB(config.DataBaseCfg)
	gameRepo := gamemysql.NewDataBase(mysqlDB)
	redisPublisher := publisher.NewPublisher(config.PublisherCfg, redisAdapter)
	redisSubscriber := subscriber.NewSubscriber(redisAdapter, config.SubscriberCfg)

	return gameservice.New(redisAdapter, gameRepo, redisPublisher, redisSubscriber, config.GameServiceCfg)
}
