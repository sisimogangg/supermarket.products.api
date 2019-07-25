package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sisimogangg/supermarket.products.api/models"
	"github.com/sisimogangg/supermarket.products.api/repository"
)

type productService struct {
	repo    repository.DataAccessLayer
	timeout time.Duration
}

// NewProductService constructs a product service instance
func NewProductService(repo repository.DataAccessLayer, timeout time.Duration) Service {
	return &productService{repo, timeout}
}

func server(productID int32) (bool, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	URL := fmt.Sprintf("http://localhost:8080/api/discount/%v", productID)
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal(res.Status)
	}

	type response struct {
		IsOnDiscount bool
		Message      string
		Status       bool
	}

	results := response{}

	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		return false, err
	}

	return results.IsOnDiscount, nil
}

func (s *productService) retrieveDiscounts(ctx context.Context, products []*models.Product) ([]*models.Product, error) {
	for _, p := range products {

		dis, err := server(p.ID)
		if err != nil {
			return nil, err
		}
		p.IsOnDiscount = dis
	}
	return products, nil
}

func (s *productService) AllProducts(ctx context.Context) ([]*models.Product, error) {
	time.Sleep(5 * time.Second)

	ps, err := s.repo.AllProducts(ctx)
	if err != nil {
		return nil, err
	}

	products, err := s.retrieveDiscounts(ctx, ps)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productService) GetProductByID(ctx context.Context, productID int32) (*models.Product, error) {
	p, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return p, nil
}
