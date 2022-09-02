package enums

type CollectionName string

const (
	USER_COLLECTION_NAME     = CollectionName("users")
	CATEGORY_COLLECTION_NAME = CollectionName("categories")
	PRODUCT_COLLECTION_NAME  = CollectionName("products")
)
