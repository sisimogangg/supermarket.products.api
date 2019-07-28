package service

import (
	"context"

	"github.com/sisimogangg/supermarket.products.api/model"
)

// Service defines the microservice contract
type Service interface {
	List(ctx context.Context) ([]*model.Product, error)
	Get(ctx context.Context, productID string) (*model.Detail, error)
}
