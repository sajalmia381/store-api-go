package service

import (
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryService interface {
	Store(payload dtos.CategoryStoreDto) (model.Category, error)
	FindAll() ([]model.Category, error)
	FindBySlug(slug string) (model.Category, error)
	UpdateBySlug(slug string, payload dtos.CategoryUpdateDto) (model.Category, error)
	DeleteBySlug(slug string) (*mongo.DeleteResult, error)
	PushProductToCategory(categoryId primitive.ObjectID, productId primitive.ObjectID) error
	RemoveProductFromCategory(categoryId primitive.ObjectID, productId primitive.ObjectID) error
	ChangeProductInCategory(categoryId primitive.ObjectID, oldProductId primitive.ObjectID, newProductId primitive.ObjectID) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func (s categoryService) Store(payload dtos.CategoryStoreDto) (model.Category, error) {
	category := model.Category{
		Name:        payload.Name,
		Description: payload.Description,
	}

	category, err := s.repo.Store(category)
	if err != nil {
		return category, err
	}
	return category, err
}

func (s categoryService) FindAll() ([]model.Category, error) {
	objects, err := s.repo.FindAll()
	return objects, err
}

func (s categoryService) FindBySlug(slug string) (model.Category, error) {
	category, err := s.repo.FindBySlug(slug)
	return category, err
}

func (s categoryService) UpdateBySlug(slug string, payload dtos.CategoryUpdateDto) (model.Category, error) {
	category, err := s.repo.UpdateBySlug(slug, payload)
	return category, err
}

func (s categoryService) DeleteBySlug(slug string) (*mongo.DeleteResult, error) {
	result, err := s.repo.DeleteBySlug(slug)
	return result, err
}

func (s categoryService) PushProductToCategory(categoryId primitive.ObjectID, productId primitive.ObjectID) error {
	err := s.repo.PushProductToCategory(categoryId, productId)
	return err
}

func (s categoryService) RemoveProductFromCategory(categoryId primitive.ObjectID, productId primitive.ObjectID) error {
	err := s.repo.RemoveProductFromCategory(categoryId, productId)
	return err
}

func (s categoryService) ChangeProductInCategory(categoryId primitive.ObjectID, oldProductId primitive.ObjectID, newProductId primitive.ObjectID) error {
	err := s.repo.ChangeProductInCategory(categoryId, oldProductId, newProductId)
	return err
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}
