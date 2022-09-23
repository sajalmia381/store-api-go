package v1

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/sajalmia381/store-api/src/api/common"
	"github.com/sajalmia381/store-api/src/config"
	"github.com/sajalmia381/store-api/src/utils"
	"github.com/sajalmia381/store-api/src/v1/api"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type cartApi struct {
	cartService service.CartService
}

// Requester Cart
func (a cartApi) ViewCart(c echo.Context) error {
	requesterId := *config.DefaultUserId
	jwtPayload, err := utils.GetRequestData(c)
	if err == nil {
		requesterId = jwtPayload.ID
	}
	cart, err := a.cartService.FindByUserId(requesterId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return common.GenerateSuccessResponse(c, nil, "User cart empty")
		}
		return common.GenerateErrorResponse(c, nil, "User cart")
	}
	return common.GenerateSuccessResponse(c, cart, "User cart")
}

func (a cartApi) UpdateCartByProducts(c echo.Context) error {
	requesterId := *config.DefaultUserId
	jwtPayload, err := utils.GetRequestData(c)
	if err == nil {
		requesterId = jwtPayload.ID
	}
	var productSpec []model.CartProductSpec
	if err := c.Bind(&productSpec); err != nil {
		log.Println("[ERROR] Cart update data bind:", err)
		return common.GenerateSuccessResponse(c, nil, "Failed to bind data")
	}
	cart, err := a.cartService.UpdateCartByProducts(requesterId, productSpec)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, cart, "Success! Cart update")
}

func (a cartApi) UpdateCartByProduct(c echo.Context) error {
	requesterId := *config.DefaultUserId
	jwtPayload, err := utils.GetRequestData(c)
	if err == nil {
		requesterId = jwtPayload.ID
	}
	var productSpec model.CartProductSpec
	if err := c.Bind(&productSpec); err != nil {
		log.Println("[ERROR] Cart update data bind:", err)
		return common.GenerateSuccessResponse(c, nil, "Failed to bind data")
	}
	cart, err := a.cartService.UpdateCartByProduct(requesterId, productSpec)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, cart, "Success! Cart update")
}
func (a cartApi) RemoveProductFromCart(c echo.Context) error {
	requesterId := *config.DefaultUserId
	jwtPayload, err := utils.GetRequestData(c)
	if err == nil {
		log.Println("[ERROR]", err)
		requesterId = jwtPayload.ID
	}
	var productIdSpec dtos.CartProductId
	if err := c.Bind(&productIdSpec); err != nil {
		log.Println("[ERROR] Cart update data bind:", err)
		return common.GenerateSuccessResponse(c, nil, "Failed to bind data")
	}
	cart, err := a.cartService.RemoveProductFromCart(requesterId, productIdSpec.ProductId)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, cart, "Success! Cart update")
}

// Cart CRUD
func (a cartApi) FindAll(c echo.Context) error {
	carts, err := a.cartService.FindAll()
	if err != nil {
		return common.GenerateErrorResponse(c, "", err.Error())
	}
	return common.GenerateSuccessResponse(c, carts, "Success! All carts list")
}

func (a cartApi) FindByUserId(c echo.Context) error {
	userId := c.Param("userId")
	_id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, "User is not valid")
	}
	cart, err := a.cartService.FindByUserId(_id)
	if err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	return common.GenerateSuccessResponse(c, cart, "User Cart")
}

func NewCartApi(service service.CartService) api.CartApi {
	return &cartApi{
		cartService: service,
	}
}
