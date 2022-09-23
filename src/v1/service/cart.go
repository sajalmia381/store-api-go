package service

import (
	"errors"

	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartService interface {
	FindAll() ([]model.Cart, error)
	FindByUserId(userId primitive.ObjectID) (model.Cart, error)
	DeleteByUserId(userId primitive.ObjectID) (*mongo.DeleteResult, error)
	// Requester Cart
	UpdateCartByProduct(userId primitive.ObjectID, payload model.CartProductSpec) (model.Cart, error)
	RemoveProductFromCart(userId primitive.ObjectID, productId string) (model.Cart, error)
	UpdateCartByProducts(userId primitive.ObjectID, payload []model.CartProductSpec) (model.Cart, error)
}

type cartService struct {
	repo repository.CartRepository
}

// Cart CRUD
func (s cartService) FindAll() ([]model.Cart, error) {
	carts, err := s.repo.FindAll()
	return carts, err
}

func (s cartService) FindByUserId(userId primitive.ObjectID) (model.Cart, error) {
	cart, err := s.repo.FindByUserId(userId)
	return cart, err
}

func (s cartService) DeleteByUserId(userId primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := s.repo.DeleteByUserId(userId)
	return result, err
}

// Requester Cart
func (s cartService) UpdateCartByProducts(userId primitive.ObjectID, productSpec []model.CartProductSpec) (model.Cart, error) {
	specMap := make(map[string]uint16)
	for _, item := range productSpec {
		_, ok := specMap[item.ProductId.Hex()]
		if ok {
			specMap[item.ProductId.Hex()] = specMap[item.ProductId.Hex()] + item.Quantity
			continue
		}
		specMap[item.ProductId.Hex()] = item.Quantity
	}
	payload := []model.CartProductSpec{}
	for key, value := range specMap {
		prodId, _ := primitive.ObjectIDFromHex(key)
		payload = append(payload, model.CartProductSpec{ProductId: prodId, Quantity: value})
	}
	cart, err := s.repo.UpdateCartByProducts(userId, payload)
	return cart, err
}

func (s cartService) UpdateCartByProduct(userId primitive.ObjectID, payload model.CartProductSpec) (model.Cart, error) {
	cart, err := s.repo.UpdateCartByProduct(userId, payload)
	return cart, err
}

func (s cartService) RemoveProductFromCart(userId primitive.ObjectID, productId string) (model.Cart, error) {
	var cart model.Cart
	_productId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return cart, errors.New("invalid ProductID id")
	}
	cart, err = s.repo.RemoveProductFromCart(userId, _productId)
	return cart, err
}

func NewCartService(cartRepo repository.CartRepository) CartService {
	return &cartService{
		repo: cartRepo,
	}
}
