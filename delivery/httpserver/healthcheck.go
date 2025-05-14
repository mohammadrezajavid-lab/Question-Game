package httpserver

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"log"
	"net/http"
)

func (hs *HttpServer) HealthCheckHandler(ctx echo.Context) error {

	db := mysql.NewDB(hs.serverConfig.DataBaseCfg)

	if err := db.MysqlConnection.Ping(); err != nil {

		log.Println(err)

		return echo.NewHTTPError(http.StatusInternalServerError, "unexpected error: ping to database server failed")
	}

	if err := db.MysqlConnection.Close(); err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "unexpected error: close database connection is failed")
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "health check OK"})
}
