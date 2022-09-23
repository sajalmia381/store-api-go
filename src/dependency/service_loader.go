package dependency

import (
	"github.com/sajalmia381/store-api/src/v1/repository"
	"github.com/sajalmia381/store-api/src/v1/service"
)

func GetTokenService() service.TokenService {
	return service.NewTokenService(repository.NewTokenRepository())
}

func GetAuthService() service.AuthService {
	return service.NewAuthService(repository.NewUserRepository(), repository.NewTokenRepository())
}

func GetJwtService() service.JwtService {
	return service.NewJwtService()
}

func GetUserService() service.UserService {
	return service.NewUserService(repository.NewUserRepository())
}

func GetCategoryService() service.CategoryService {
	return service.NewCategoryService(repository.NewCategoryRepository())
}

func GetProductService() service.ProductService {
	return service.NewProductService(repository.NewProductRepository(), repository.NewCategoryRepository())
}

func GetCartService() service.CartService {
	return service.NewCartService(repository.NewCartRepository())
}
