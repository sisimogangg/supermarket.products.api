package product

import (
	"context"

	"github.com/sisimogangg/supermarket.products.api/models"
)

// Service defines the microservice contract
type Service interface {
	AllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, productID int) (*models.Product, error)
}
