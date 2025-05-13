package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret_key"),
	}))

	e.GET("/", func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {

			return errors.New("JWT token missing or invalid")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {

			return errors.New("failed to cast claims as jwt.MapClaims")
		}

		return c.JSON(http.StatusOK, claims)
	})

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
