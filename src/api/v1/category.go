package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/sajalmia381/store-api/src/api/common"
	"github.com/sajalmia381/store-api/src/utils"
	"github.com/sajalmia381/store-api/src/v1/api"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/service"
)

type categoryApi struct {
	categoryService service.CategoryService
}

func (cat categoryApi) Store(c echo.Context) error {
	var formData dtos.CategoryStoreDto
	if err := c.Bind(&formData); err != nil {
		return common.GenerateErrorResponse(c, err, "Failed to bind data")
	}
	if err := formData.Validate(); err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	isSuperAdmin := utils.IsSuperAdmin(c)
	var (
		category model.Category
		err      error
	)
	if !isSuperAdmin {
		category = cat.categoryService.FakeStore(formData)
	} else {
		category, err = cat.categoryService.Store(formData)
	}
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, category, "Success! Category created")
}

func (cat categoryApi) FindAll(c echo.Context) error {
	categories, err := cat.categoryService.FindAll()
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, categories, "Success! Category list")
}

func (cat categoryApi) FindBySlug(c echo.Context) error {
	slug := c.Param("slug")
	category, err := cat.categoryService.FindBySlug(slug)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, category, "Success! Category description")
}

func (cat categoryApi) UpdateBySlug(c echo.Context) error {
	slug := c.Param("slug")
	var formData dtos.CategoryUpdateDto
	if err := c.Bind(&formData); err != nil {
		return common.GenerateErrorResponse(c, formData, "Failed to bind data")
	}
	isSuperAdmin := utils.IsSuperAdmin(c)
	var (
		category model.Category
		err      error
	)
	if !isSuperAdmin {
		category, err = cat.categoryService.FakeUpdateBySlug(slug, formData)
	} else {
		category, err = cat.categoryService.UpdateBySlug(slug, formData)
	}
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, category, "Success! category updated")
}

func (cat categoryApi) DeleteBySlug(c echo.Context) error {
	slug := c.Param("slug")
	isSuperAdmin := utils.IsSuperAdmin(c)
	var (
		err error
	)
	if !isSuperAdmin {
		_, err = cat.categoryService.FindBySlug(slug)
	} else {
		_, err = cat.categoryService.DeleteBySlug(slug)
	}
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, nil, "Success! Category deleted")
}

func NewCategoryApi(categoryService service.CategoryService) api.CategoryApi {
	return &categoryApi{
		categoryService: categoryService,
	}
}
