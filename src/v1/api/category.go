package api

import "github.com/labstack/echo/v4"

type CategoryApi interface {
	Store(c echo.Context) error
	FindAll(c echo.Context) error
	FindBySlug(c echo.Context) error
	UpdateBySlug(c echo.Context) error
	DeleteBySlug(c echo.Context) error
}
