package api

import "github.com/labstack/echo/v4"

type CartApi interface {
	// Cart CRUD
	FindAll(c echo.Context) error
	// FindById(c echo.Context) error // TODO: Implement
	FindByUserId(c echo.Context) error
	// DeleteById(c echo.Context) error // TODO: Implement

	// Requester Cart!!! By default all requester request as anonymous@gmail.com user
	UpdateCartByProducts(c echo.Context) error
	ViewCart(c echo.Context) error
	UpdateCartByProduct(c echo.Context) error
	RemoveProductFromCart(c echo.Context) error
}
