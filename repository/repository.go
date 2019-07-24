package repository

import (
	"context"

	"github.com/sisimogangg/supermarket.products.api/models"
)

// DataAccessLayer defines expected repository behavour
type DataAccessLayer interface {
	AllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, productID int) (*models.Product, error)
}
