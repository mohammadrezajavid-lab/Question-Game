package main

import (
	"context"
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/wsconfig"
	"golang.project/go-fundamentals/gameapp/gateway/websocket"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/infomessage"
	"os"
	"os/signal"
	"sync"
)

func main() {

	var host string
	var port int
	flag.StringVar(&host, "host", "", "webSocket http server host")
	flag.IntVar(&port, "port", 0, "webSocket http server port")
	flag.Parse()

	allConfig := wsconfig.NewConfig().LoadConfig(host, port)

	logger.InitLogger(allConfig.LoggerCfg)

	logger.Info(fmt.Sprintf("webSocket config: %v", allConfig))

	ws, sErr := websocket.NewWebSocket(allConfig.WSCfg, allConfig.JwtCfg, allConfig.GrpcPresenceClientCfg)
	if sErr != nil {
		logger.Fatal(sErr, "")
	}
	go ws.ServeWs()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-ctx.Done()

	logger.Info(infomessage.InfoMsgShuttingDownGracefully)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), allConfig.WSCfg.GracefullyShutdownTimeout)
	defer cancel()

	var shutdownWG sync.WaitGroup

	shutdownWG.Add(1)
	go func() {
		defer shutdownWG.Done()
		if err := ws.Shutdown(shutdownCtx); err != nil {
			logger.Error(err, "failed to gracefully shutdown WebSocket gateway")
		}
	}()

	shutdownWG.Wait()
}
