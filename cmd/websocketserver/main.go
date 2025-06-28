package main

import (
	"flag"
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/wsconfig"
	"golang.project/go-fundamentals/gameapp/gateway/websocket"
	"golang.project/go-fundamentals/gameapp/logger"
)

func main() {

	var host string
	var port int
	flag.StringVar(&host, "host", "", "webSocket http server host")
	flag.IntVar(&port, "port", 0, "webSocket http server port")
	flag.Parse()

	wsCfg := wsconfig.NewConfig().LoadConfig(host, port)

	logger.InitLogger(wsCfg.LoggerCfg)

	logger.Info(fmt.Sprintf("webSocket config: %v", wsCfg))

	ws := websocket.NewWebSocket(wsCfg.WSCfg, wsCfg.JwtCfg)
	ws.ServeWs()
}
