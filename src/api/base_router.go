package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	v1 "github.com/sajalmia381/store-api/src/api/v1"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Routes(e *echo.Echo) {
	e.GET("/", index)
	e.GET("/health", health)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	v1.Routes(e.Group("/api/v1"))
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to Store-Api")
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "I am alive...")
}
