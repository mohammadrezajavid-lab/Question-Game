package httpserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"log"
	"net/http"
	"time"
)

func (s *Server) HealthCheckHandler(c echo.Context) error {

	ctx := c.Request().Context()

	select {
	case <-time.After(9 * time.Second):
		return c.JSON(http.StatusGatewayTimeout, echo.Map{
			"message": "server timeout",
		})
	case <-ctx.Done():
		err := ctx.Err()
		fmt.Println("server: ", err)
	default:

		db := mysql.NewDB(s.config.DataBaseCfg)

		errorFields := make([]error, 0)

		if err := db.MysqlConnection.Ping(); err != nil {

			log.Println(err)

			errorFields = append(errorFields, fmt.Errorf("unexpected error: ping to database server failed"))
		}
		if err := db.MysqlConnection.Close(); err != nil {

			errorFields = append(errorFields, fmt.Errorf("unexpected error: close database connection is failed"))
		}

		redisAdapter := redis.New(s.config.RedisCfg)
		rdb := redisAdapter.GetClient()
		if _, pErr := rdb.Ping(ctx).Result(); pErr != nil {
			log.Println(pErr)

			errorFields = append(errorFields, fmt.Errorf("unexpected error: ping to redis server failed"))
		}
		if err := rdb.Close(); err != nil {

			errorFields = append(errorFields, fmt.Errorf("unexpected error: close redis connection is failed"))
		}

		if len(errorFields) != 0 {

			return c.JSON(http.StatusInternalServerError, echo.Map{
				"errors": strErrors(errorFields),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{"message": "health check OK"})
	}

	return nil
}

func strErrors(errs []error) []string {
	stringErrors := make([]string, 0)
	for _, err := range errs {
		stringErrors = append(stringErrors, err.Error())
	}

	return stringErrors
}
