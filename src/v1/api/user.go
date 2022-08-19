package api

import "github.com/labstack/echo/v4"

type User interface {
	Store(c echo.Context) error
	FindAll(c echo.Context) error
	FindById(c echo.Context) error
	UpdateById(c echo.Context) error
	DeleteById(c echo.Context) error
}
