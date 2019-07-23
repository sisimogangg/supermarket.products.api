package repository

import (
	"context"

	"github.com/sisimogangg/supermarket.products.api/models"
	"github.com/sisimogangg/supermarket.products.api/product"
)

type firebaseRepo struct{}

var products = [...]models.Product{
	models.Product{
		ID:       1,
		ImageURL: "https://images.unsplash.com/photo-1478004521390-655bd10c9f43?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
		Name:     "Apple",
		Price: struct {
			Symbol   rune
			Currency string
			Amount   float32
		}{
			Symbol:   'R',
			Currency: "RSA",
			Amount:   2.00,
		},
		Promotion: "",
	},
	models.Product{
		ID:       2,
		ImageURL: "https://images.unsplash.com/photo-1528825871115-3581a5387919?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=658&q=80",
		Name:     "Banana",
		Price: struct {
			Symbol   rune
			Currency string
			Amount   float32
		}{
			Symbol:   'R',
			Currency: "RSA",
			Amount:   3.00,
		},
		Promotion: "",
	},
	models.Product{
		ID:       3,
		ImageURL: "https://images.unsplash.com/photo-1560769680-ba2f3767c785?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjExMDk0fQ&auto=format&fit=crop&w=700&q=80",
		Name:     "Coconut",
		Price: struct {
			Symbol   rune
			Currency string
			Amount   float32
		}{
			Symbol:   'R',
			Currency: "RSA",
			Amount:   4.00,
		},
		Promotion: "",
	},
}

// NewFirebaseRepo defines a constructor for firebaserepo
func NewFirebaseRepo() product.Repository {
	return &firebaseRepo{}
}

//AllProducts returns all products
func (f *firebaseRepo) AllProducts(ctx context.Context) ([]*models.Product, error) {
	ps := make([]*models.Product, 0)
	for _, p := range products {
		ps = append(ps, &p)
	}
	return ps, nil
}

//GetProductByID returns product given its ID
func (f *firebaseRepo) GetProductByID(ctx context.Context, productID int) (*models.Product, error) {
	product := models.Product{}
	for _, p := range products {
		if p.ID == productID {
			product = p
		}
	}
	if product.ID == 0 {
		return nil, nil
	}
	return &product, nil
}
