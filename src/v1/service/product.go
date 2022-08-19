package service

import (
	"errors"
	"log"

	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductService interface {
	Store(payload dtos.ProductStoreDto) (model.Product, error)
	FindAll() ([]dtos.ProductResponseDto, error)
	FindBySlug(slug string) (model.Product, error)
	UpdateBySlug(slug string, payload dtos.ProductUpdateDto) (model.Product, error)
	DeleteBySlug(slug string) (*mongo.DeleteResult, error)
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

func (p productService) FindAll() ([]dtos.ProductResponseDto, error) {
	products, err := p.repo.FindAll()
	return products, err
}

func (p productService) FindBySlug(slug string) (model.Product, error) {
	product, err := p.repo.FindBySlug(slug)
	return product, err
}

func (p productService) UpdateBySlug(slug string, payload dtos.ProductUpdateDto) (model.Product, error) {
	product, err := p.FindBySlug(slug)
	if err != nil {
		return product, err
	}

	log.Println(payload.Category, product.Category)

	isCategoryChange := payload.Category.Hex() != "" && payload.Category.Hex() != product.Category.Hex()
	log.Println(isCategoryChange)
	newProduct, err := p.repo.UpdateOne(product, payload)
	if err == nil {
		log.Println(payload.Category, product.Category)
		log.Println(payload.Category.Hex() != "" && payload.Category.Hex() != product.Category.Hex())
		if isCategoryChange {
			go func() {
				err := p.categoryRepo.ChangeProductInCategory(*product.Category, *payload.Category, product.ID)
				if err != nil {
					log.Println("[ERROR update product category change]:", err.Error())
				}
			}()
		}
	}
	return newProduct, err
}

func (p productService) DeleteBySlug(slug string) (*mongo.DeleteResult, error) {
	product, err := p.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	result, err := p.repo.DeleteBySlug(slug)
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
	return result, err
}

func NewProductService(repo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}
