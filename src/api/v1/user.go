package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sajalmia381/store-api/src/api/common"
	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/utils"
	"github.com/sajalmia381/store-api/src/v1/api"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
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

	isSuperAdmin := utils.IsSuperAdmin(c)
	var (
		user model.User
		err  error
	)
	if !isSuperAdmin {
		user, err = u.userService.FakeStore(formData)
	} else {
		user, err = u.userService.Store(formData)
	}

	// To super admin
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, user, "Success! User created", &common.ResponseOption{
		HttpCode: http.StatusCreated,
	})
}

func (u userApi) FindAll(c echo.Context) error {
	// log.Println("req", c.Request())
	isSuperAdmin := utils.IsSuperAdmin(c)
	var (
		objects []model.User
		err     error
	)
	query := dtos.UserQuery{}
	if isSuperAdmin {
		objects, err = u.userService.FindAll(query)
	} else {
		query.Role = string(enums.ROLE_CUSTOMER)
		objects, err = u.userService.FindAll(query)
	}
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
	isSuperAdmin := utils.IsSuperAdmin(c)
	var (
		user model.User
		err  error
	)
	if !isSuperAdmin {
		user, err = u.userService.FakeUpdateById(id, formData)
	} else {
		user, err = u.userService.UpdateById(id, formData)
	}
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}

	return common.GenerateSuccessResponse(c, user, "Success! User updated")
}

func (u userApi) DeleteById(c echo.Context) error {
	id := c.Param("id")
	isSuperAdmin := utils.IsSuperAdmin(c)
	var err error
	if !isSuperAdmin {
		_, err = u.userService.FindById(id)
	} else {
		err = u.userService.DeleteById(id)
	}
	if err != nil {
		return common.GenerateErrorResponse(c, "Failed to delete user", err.Error())
	}
	return common.GenerateSuccessResponse(c, nil, "Success! User deleted")
}

func NewUserApi(userService service.UserService) api.User {
	return &userApi{
		userService: userService,
	}
}
