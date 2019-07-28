package repository

import (
	"context"

	"github.com/sisimogangg/supermarket.products.api/model"
)

// DataAccessLayer defines expected repository behavour
type DataAccessLayer interface {
	List(ctx context.Context) ([]*model.Product, error)
	Get(ctx context.Context, productID string) (*model.Detail, error)
}
