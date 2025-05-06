package config

import (
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"golang.project/go-fundamentals/gameapp/service/authorize"
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
	AuthCfg       authorize.Config
	DataBaseCfg   mysql.Config
}

func NewConfig(httpServerCfg HttpServerCfg, authCfg authorize.Config, dataBaseCfg mysql.Config) Config {
	return Config{
		HttpServerCfg: httpServerCfg,
		AuthCfg:       authCfg,
		DataBaseCfg:   dataBaseCfg,
	}
}
