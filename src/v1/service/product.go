package service

import (
	"errors"
	"log"
	"time"

	"github.com/sajalmia381/store-api/src/api/common"
	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/utils"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	Store(payload dtos.ProductStoreDto) (model.Product, error)
	FindAll(queryParams dtos.ProductQueryParams) ([]dtos.ProductResponseDto, common.MetaData, error)
	FindBySlug(slug string) (model.Product, error)
	UpdateBySlug(slug string, payload dtos.ProductUpdateDto) (model.Product, error)
	DeleteBySlug(slug string) (model.Product, error)
	// Fake
	FakeStore(payload dtos.ProductStoreDto) (model.Product, error)
	FakeUpdateBySlug(slug string, payload dtos.ProductUpdateDto) (model.Product, error)
}

type productService struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func (p productService) Store(payload dtos.ProductStoreDto) (model.Product, error) {
	product := model.Product{
		Title:       payload.Title,
		Description: *payload.Description,
		Price:       *payload.Price,
	}
	catId, err := primitive.ObjectIDFromHex(payload.Category)
	if err != nil {
		return product, errors.New("category id is not valid")
	}
	product.Category = &catId
	// No dep
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now().UTC()
	product.UpdatedAt = time.Now().UTC()
	product.Active = true
	product, err = p.repo.Store(product)
	if err != nil {
		return product, err
	}
	if product.Category.Hex() != "" {
		go func() {
			err := p.categoryRepo.PushProductToCategory(*product.Category, product.ID)
			if err != nil {
				log.Println("[ERROR] add product in category:", err.Error())
			}
		}()
	}
	return product, nil
}

func (p productService) FindAll(queryParams dtos.ProductQueryParams) ([]dtos.ProductResponseDto, common.MetaData, error) {
	products, metaData, err := p.repo.FindAll(queryParams)
	return products, metaData, err
}

func (p productService) FindBySlug(slug string) (model.Product, error) {
	product, err := p.repo.FindBySlug(slug)
	return product, err
}

func (p productService) UpdateBySlug(slug string, formData dtos.ProductUpdateDto) (model.Product, error) {
	payload := primitive.M{}

	if formData.Title != "" {
		payload["title"] = formData.Title
	}
	if formData.Price != nil {
		payload["price"] = &formData.Price
	}
	if formData.Description != "" {
		payload["description"] = formData.Description
	}
	if formData.Category != nil {
		payload["category"] = formData.Category
	}
	if formData.UpdateSlug && formData.Title != "" {
		payload["slug"] = utils.GenerateUniqueSlug(formData.Title, string(enums.PRODUCT_COLLECTION_NAME), slug)
	}

	product, err := p.repo.UpdateBySlug(slug, payload)
	isCategoryChange := formData.Category.Hex() != "" && formData.Category.Hex() != product.Category.Hex()
	log.Println(isCategoryChange)
	if err == nil {
		log.Println(formData.Category, formData.Category)
		log.Println(formData.Category.Hex() != "" && formData.Category.Hex() != product.Category.Hex())
		// TODO: Fix product change to push category
		if isCategoryChange {
			go func() {
				err := p.categoryRepo.ChangeProductInCategory(*product.Category, *formData.Category, product.ID)
				if err != nil {
					log.Println("[ERROR update product category change]:", err.Error())
				}
			}()
		}
	}
	return product, err
}

func (p productService) DeleteBySlug(slug string) (model.Product, error) {
	product, err := p.repo.DeleteBySlug(slug)
	if err == nil {
		if product.Category.Hex() != "" {
			go func() {
				err := p.categoryRepo.RemoveProductFromCategory(*product.Category, product.ID)
				if err != nil {
					log.Println("[ERROR] delete product from category:", err.Error())
				}
			}()
		}
	}
	return product, err
}

func (p productService) FakeStore(payload dtos.ProductStoreDto) (model.Product, error) {
	product := model.Product{
		Title:       payload.Title,
		Description: *payload.Description,
		Price:       *payload.Price,
	}
	catId, err := primitive.ObjectIDFromHex(payload.Category)
	if err != nil {
		return product, errors.New("category id is not valid")
	}
	product.Category = &catId
	// No dep
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now().UTC()
	product.UpdatedAt = time.Now().UTC()
	product.Active = true
	product.Slug = utils.GenerateFakeUniqueSlug(product.Title, false)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (p productService) FakeUpdateBySlug(slug string, payload dtos.ProductUpdateDto) (model.Product, error) {
	product, err := p.FindBySlug(slug)
	if err != nil {
		return product, err
	}
	if payload.Title != "" {
		product.Title = payload.Title
	}
	if payload.Price != nil {
		product.Price = *payload.Price
	}
	if payload.Description != "" {
		product.Description = payload.Description
	}
	if payload.Category != nil {
		product.Category = payload.Category
	}
	if payload.UpdateSlug && payload.Title != "" {
		product.Slug = utils.GenerateFakeUniqueSlug(payload.Title, false)
	}
	return product, err
}

func NewProductService(repo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}
