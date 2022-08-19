package api

import (
	"github.com/labstack/echo/v4"
)

type AuthApi interface {
	Login(c echo.Context) error
	Register(c echo.Context) error
	RefreshToken(c echo.Context) error
}
