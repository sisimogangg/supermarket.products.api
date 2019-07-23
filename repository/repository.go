package repository

import (
	"context"

	"github.com/sisimogangg/supermarket.products.api/models"
)

// Repository defines expected repository behavour
type Repository interface {
	AllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, productID int) (*models.Product, error)
}
