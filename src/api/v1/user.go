package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/sajalmia381/store-api/src/api/common"
	"github.com/sajalmia381/store-api/src/v1/api"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/service"
)

type userApi struct {
	userService service.UserService
}

func (u userApi) Store(c echo.Context) error {
	var formData dtos.UserRegisterDTO
	if err := formData.Validate(); err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	if err := c.Bind(&formData); err != nil {
		return common.GenerateErrorResponse(c, err.Error(), "Failed to bind data!")
	}
	user, err := u.userService.Store(formData)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, user, "Success! User created", &common.ResponseOption{
		HttpCode: 201,
	})
}

func (u userApi) FindAll(c echo.Context) error {
	objects, err := u.userService.FindAll()
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, objects, "Success! User list")
}

func (u userApi) FindById(c echo.Context) error {
	id := c.Param("id")
	user, err := u.userService.FindById(id)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, user, "Success! User description")
}

func (u userApi) UpdateById(c echo.Context) error {
	id := c.Param("id")
	var formData dtos.UserUpdateDto
	if err := formData.Validate(); err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	if err := c.Bind(&formData); err != nil {
		return common.GenerateErrorResponse(c, err.Error(), "Failed to bind data!")
	}
	user, err := u.userService.UpdateById(id, formData)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}

	return common.GenerateSuccessResponse(c, user, "Success! User updated")
}

func (u userApi) DeleteById(c echo.Context) error {
	id := c.Param("id")
	err := u.userService.DeleteById(id)
	if err != nil {
		return common.GenerateErrorResponse(c, "Failed to delete user", err.Error())
	}
	return common.GenerateSuccessResponse(c, nil, "Success! User delete")
}

func NewUserApi(userService service.UserService) api.User {
	return &userApi{
		userService: userService,
	}
}
