package service

import (
	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/utils"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/repository"
	"go.mongodb.org/mongo-driver/bson"
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
	// Fake Action
	FakeStore(payload dtos.CategoryStoreDto) model.Category
	FakeUpdateBySlug(slug string, payload dtos.CategoryUpdateDto) (model.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func (s categoryService) Store(payload dtos.CategoryStoreDto) (model.Category, error) {
	category := model.Category{
		ID:          primitive.NewObjectID(),
		Name:        payload.Name,
		Description: payload.Description,
		Products:    []primitive.ObjectID{},
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

func (s categoryService) UpdateBySlug(slug string, formData dtos.CategoryUpdateDto) (model.Category, error) {

	payload := bson.M{}
	if formData.Name != "" {
		payload["name"] = formData.Name
	}
	if formData.Description != "" {
		payload["description"] = formData.Description
	}
	if formData.UpdateSlug && formData.Name != "" {
		payload["slug"] = utils.GenerateUniqueSlug(formData.Name, string(enums.CATEGORY_COLLECTION_NAME), slug)
	}
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

// Fake
func (s categoryService) FakeStore(payload dtos.CategoryStoreDto) model.Category {
	category := model.Category{
		ID:          primitive.NewObjectID(),
		Name:        payload.Name,
		Description: payload.Description,
		Products:    []primitive.ObjectID{},
		Slug:        utils.GenerateFakeUniqueSlug(payload.Name, false),
	}
	return category
}

func (s categoryService) FakeUpdateBySlug(slug string, payload dtos.CategoryUpdateDto) (model.Category, error) {
	category, err := s.FindBySlug(slug)
	if err != nil {
		return category, err
	}
	if payload.Name != "" {
		category.Name = payload.Name
	}
	if payload.Description != "" {
		category.Description = payload.Description
	}
	if payload.UpdateSlug && payload.Name != "" {
		category.Slug = utils.GenerateFakeUniqueSlug(payload.Name, false)
	}
	return category, nil
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}
