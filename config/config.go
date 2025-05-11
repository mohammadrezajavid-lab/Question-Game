package config

import (
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/service/auth"
)

type HttpServerCfg struct {
	Host string
	Port int
}

func NewHttpServerCfg(host string, port int) HttpServerCfg {
	return HttpServerCfg{Host: host, Port: port}
}

type Config struct {
	HttpServerCfg HttpServerCfg
	AuthCfg       auth.Config
	DataBaseCfg   mysql.Config
}

func NewConfig(httpServerCfg HttpServerCfg, authCfg auth.Config, dataBaseCfg mysql.Config) Config {
	return Config{
		HttpServerCfg: httpServerCfg,
		AuthCfg:       authCfg,
		DataBaseCfg:   dataBaseCfg,
	}
}
