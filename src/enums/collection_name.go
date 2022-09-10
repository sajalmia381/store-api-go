package enums

type CollectionName string

const (
	TOKEN_COLLECTION_NAME    = CollectionName("tokens")
	USER_COLLECTION_NAME     = CollectionName("users")
	CATEGORY_COLLECTION_NAME = CollectionName("categories")
	PRODUCT_COLLECTION_NAME  = CollectionName("products")
	CART_COLLECTION_NAME     = CollectionName("carts")
	ORDER_COLLECTION_NAME    = CollectionName("orders")
)

var COLLECTION_NAMES = []string{
	string(TOKEN_COLLECTION_NAME),
	string(USER_COLLECTION_NAME),
	string(CATEGORY_COLLECTION_NAME),
	string(PRODUCT_COLLECTION_NAME),
	string(CART_COLLECTION_NAME),
	string(ORDER_COLLECTION_NAME),
}
