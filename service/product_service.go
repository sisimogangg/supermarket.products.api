package service

import (
	"context"
	"time"

	"github.com/sisimogangg/supermarket.products.api/models"
	"github.com/sisimogangg/supermarket.products.api/repository"
)

type productService struct {
	repo    repository.Repository
	timeout time.Duration
}

// NewProductService constructs a product service instance
func NewProductService(repo repository.Repository, timeout time.Duration) Service {
	return &productService{repo, timeout}
}

func (s *productService) AllProducts(ctx context.Context) ([]*models.Product, error) {
	ps, err := s.repo.AllProducts(ctx)
	if err != nil {
		return nil, err
	}

	return ps, nil
}
func (s *productService) GetProductByID(ctx context.Context, productID int) (*models.Product, error) {
	p, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return p, nil
}
