package config

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	IntVariables()
	isConnected := InitDBConnection()
	if isConnected {
		log.Println("[INFO] Database connected...")
	}
	go DBHealthChecker()

	echoInstance := echo.New()

	// Configuring Middleware Logger
	echoInstance.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		// Skipping logging for health checking api
		Skipper: func(c echo.Context) bool {
			return c.Request().RequestURI == "/health"
		},
		Format: "[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}, latency=${latency_human} remote_ip=${remote_ip}\n",
	}))

	echoInstance.Use(middleware.Recover())

	echoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	return echoInstance
}
