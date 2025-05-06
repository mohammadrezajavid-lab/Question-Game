package httpserver

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"log"
	"net/http"
)

func (hs *HttpServer) HealthCheckHandler(ctx echo.Context) error {

	db := mysql.NewDB(hs.Config.DataBaseCfg)
	if err := db.MysqlConnection.Ping(); err != nil {

		log.Println(err)

		return ctx.JSON(http.StatusInternalServerError, "unexpected error: ping to database server failed")
	}

	if err := db.MysqlConnection.Close(); err != nil {

		return ctx.JSON(http.StatusInternalServerError, "unexpected error: close database connection is failed")
	}

	return ctx.JSON(http.StatusOK, "health check OK")
}
