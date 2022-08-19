package dtos

import "errors"

type CategoryUpdateDto struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Slug        string `json:"slug" bson:"slug"`
	UpdateSlug  bool   `json:"updateSlug" bson:"updateSlug"`
}

type CategoryStoreDto struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

func (c *CategoryStoreDto) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
