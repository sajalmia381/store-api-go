package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/sajalmia381/store-api/src/api/common"
	"github.com/sajalmia381/store-api/src/v1/api"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/service"
)

type productApi struct {
	productService service.ProductService
}

func (p productApi) Store(c echo.Context) error {
	var formData dtos.ProductStoreDto
	if err := c.Bind(&formData); err != nil {
		return common.GenerateErrorResponse(c, nil, "Failed to bind data")
	}
	if err := formData.Validate(); err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	product, err := p.productService.Store(formData)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, product, "Success! Product description")
}

func (p productApi) FindAll(c echo.Context) error {
	products, err := p.productService.FindAll()
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, products, "Success! Product list")
}

func (p productApi) FindBySlug(c echo.Context) error {
	slug := c.Param("slug")
	product, err := p.productService.FindBySlug(slug)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, product, "Success! Product description")
}

func (p productApi) UpdateBySlug(c echo.Context) error {
	slug := c.Param("slug")
	var formData dtos.ProductUpdateDto
	if err := c.Bind(&formData); err != nil {
		return common.GenerateErrorResponse(c, nil, "Failed to bind data")
	}
	product, err := p.productService.UpdateBySlug(slug, formData)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, product, "Success! Product description")
}

func (p productApi) DeleteBySlug(c echo.Context) error {
	slug := c.Param("slug")
	_, err := p.productService.DeleteBySlug(slug)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, nil, "Success! Product deleted")
}

func NewProductApi(productService service.ProductService, categoryService service.CategoryService) api.ProductApi {
	return &productApi{
		productService: productService,
	}
}
