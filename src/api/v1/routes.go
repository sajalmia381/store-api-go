package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sajalmia381/store-api/src/custom_middleware"
	"github.com/sajalmia381/store-api/src/dependency"
)

func Routes(e *echo.Echo) {
	authRoutes(e.Group("/auth"))
	userRoutes(e.Group("/users"))
	categoryRoutes(e.Group("/categories"))
	productRoutes(e.Group("/products"))
	cartCrudRoutes(e.Group("/carts"))
	cartRequesterRoutes(e.Group("/cart"))
}

func authRoutes(g *echo.Group) {
	newAuthApi := NewAuthApi(dependency.GetAuthService(), dependency.GetJwtService(), dependency.GetTokenService())
	g.POST("/login", newAuthApi.Login)
	g.POST("/register", newAuthApi.Register)
	g.POST("/refresh", newAuthApi.RefreshToken)
}

func userRoutes(g *echo.Group) {
	newUserApi := NewUserApi(dependency.GetUserService())
	g.Use(middleware.JWTWithConfig(custom_middleware.AttachUserMiddlewareConfig()))
	g.GET("", newUserApi.FindAll)
	g.POST("", newUserApi.Store)
	g.GET("/:id", newUserApi.FindById)
	g.PUT("/:id", newUserApi.UpdateById)
	g.DELETE("/:id", newUserApi.DeleteById)
}

func categoryRoutes(g *echo.Group) {
	newCategoryApi := NewCategoryApi(dependency.GetCategoryService())
	g.Use(middleware.JWTWithConfig(custom_middleware.AttachUserMiddlewareConfig()))
	g.GET("/all", newCategoryApi.FindAll)
	g.GET("/:slug", newCategoryApi.FindBySlug)
	g.PUT("/:slug", newCategoryApi.UpdateBySlug)
	g.DELETE("/:slug", newCategoryApi.DeleteBySlug)
}

func productRoutes(g *echo.Group) {
	newProductApi := NewProductApi(dependency.GetProductService(), dependency.GetCategoryService())
	g.Use(middleware.JWTWithConfig(custom_middleware.AttachUserMiddlewareConfig()))
	g.GET("", newProductApi.FindAll)
	g.POST("", newProductApi.Store)
	g.GET("/:slug", newProductApi.FindBySlug)
	g.PUT("/:slug", newProductApi.UpdateBySlug)
	g.DELETE("/:slug", newProductApi.DeleteBySlug)
}

func cartRequesterRoutes(g *echo.Group) {
	newCartApi := NewCartApi(dependency.GetCartService())
	g.Use(middleware.JWTWithConfig(custom_middleware.AttachUserMiddlewareConfig()))
	g.GET("", newCartApi.ViewCart)                     // Request User cart
	g.PUT("", newCartApi.UpdateCartByProducts)         // Request User update Cart with bulk product
	g.PUT("/add", newCartApi.UpdateCartByProduct)      // Request User Add To Cart with single product
	g.PUT("/remove", newCartApi.RemoveProductFromCart) // Request User Remove From Cart with single product
}

func cartCrudRoutes(g *echo.Group) {
	newCartApi := NewCartApi(dependency.GetCartService())
	g.GET("", newCartApi.FindAll) // All Carts
	g.GET("/:userId", newCartApi.FindByUserId)
}
