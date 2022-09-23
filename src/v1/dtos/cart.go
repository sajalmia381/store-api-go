package dtos

import "github.com/sajalmia381/store-api/src/v1/model"

type CartProductSpecRes struct {
	Product  model.Product `json:"product" bson:"product"`
	Quantity uint16        `json:"quantity" bson:"quantity"`
	Total    uint32        `json:"total" bson:"total"`
}

type CartProductId struct {
	ProductId string `json:"productId"`
}
