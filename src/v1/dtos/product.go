package dtos

import (
	"errors"
	"time"

	"github.com/sajalmia381/store-api/src/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductStoreDto struct {
	Category    string  `json:"category" bson:"category"`
	Title       string  `json:"title" bson:"title"`
	Price       *int    `json:"price" bson:"price"`
	Description *string `json:"description" bson:"description"`
}

func (p ProductStoreDto) Validate() error {
	if p.Title == "" {
		return errors.New("title is required")
	}
	if p.Price == nil {
		return errors.New("price is required")
	}

	if p.Category == "" {
		return errors.New("category is required")
	}
	return nil
}

type ProductUpdateDto struct {
	Category    *primitive.ObjectID `json:"category" bson:"category"`
	Title       string              `json:"title" bson:"title"`
	Price       *int                `json:"price" bson:"price"`
	Description string              `json:"description" bson:"description"`
	Slug        string              `json:"slug" bson:"slug"`
	UpdateSlug  bool                `json:"updateSlug" bson:"updateSlug"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
}

// All Product

type category struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Slug string             `json:"slug" bson:"slug"`
}

type user struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Number    *int               `json:"number" bson:"number"`
	Status    bool               `json:"status" bson:"status"`
	Role      enums.Role         `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdateAt  time.Time          `json:"updateAt" bson:"updateAt"`
}

type ProductResponseDto struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedBy   *user              `json:"createdBy" bson:"createdBy"`
	Category    *category          `json:"category" bson:"category"`
	Title       string             `json:"title" bson:"title"`
	Slug        string             `json:"slug" bson:"slug"`
	Price       *int               `json:"price" bson:"price"`
	Description *string            `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	Active      bool               `json:"active" bson:"active"`
}
