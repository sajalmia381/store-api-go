package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sajalmia381/store-api/src/custom_middleware"
	"github.com/sajalmia381/store-api/src/dependency"
)

func Routes(g *echo.Group) {
	authRoutes(g.Group("/auth"))
	userRoutes(g.Group("/users"))
	categoryRoutes(g.Group("/categories"))
	productRoutes(g.Group("/products"))
}

func authRoutes(g *echo.Group) {
	newAuthApi := NewAuthApi(dependency.GetAuthService(), dependency.GetJwtService())
	g.POST("/login", newAuthApi.Login)
	g.POST("/register", newAuthApi.Register)
	g.POST("/refresh", newAuthApi.RefreshToken)
}

func userRoutes(g *echo.Group) {
	newUserApi := NewUserApi(dependency.GetUserService())
	g.Use(middleware.JWTWithConfig(custom_middleware.AuthMiddlewareConfig()))
	g.GET("", newUserApi.FindAll)
	g.POST("", newUserApi.Store)
	g.GET("/:id", newUserApi.FindById)
	g.PUT("/:id", newUserApi.UpdateById)
	g.DELETE("/:id", newUserApi.DeleteById)
}

func categoryRoutes(g *echo.Group) {
	newCategoryApi := NewCategoryApi(dependency.GetCategoryService())
	g.Use(middleware.JWTWithConfig(custom_middleware.AuthMiddlewareConfig()))
	g.GET("", newCategoryApi.FindAll)
	g.POST("", newCategoryApi.Store)
	g.GET("/:slug", newCategoryApi.FindBySlug)
	g.PUT("/:slug", newCategoryApi.UpdateBySlug)
	g.DELETE("/:slug", newCategoryApi.DeleteBySlug)
}

func productRoutes(g *echo.Group) {
	newProductApi := NewProductApi(dependency.GetProductService(), dependency.GetCategoryService())
	g.Use(middleware.JWTWithConfig(custom_middleware.AuthMiddlewareConfig()))
	g.GET("", newProductApi.FindAll)
	g.POST("", newProductApi.Store)
	g.GET("/:slug", newProductApi.FindBySlug)
	g.PUT("/:slug", newProductApi.UpdateBySlug)
	g.DELETE("/:slug", newProductApi.DeleteBySlug)
}
