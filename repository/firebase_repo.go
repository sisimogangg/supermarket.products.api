package repository

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/sisimogangg/supermarket.products.api/model"
)

type firebaseRepo struct {
	fb *firebase.App
}

// NewFirebaseRepo defines a constructor for firebaserepo
func NewFirebaseRepo(app *firebase.App) DataAccessLayer {
	return &firebaseRepo{app}
}

//List returns all products
func (f *firebaseRepo) List(ctx context.Context) ([]*model.Product, error) {
	client, err := f.fb.Database(ctx)
	productsRef := client.NewRef("products")
	if err != nil {
		return nil, err // internal server error
	}

	products := []*model.Product{}
	var rawResult map[string]model.Product

	err = productsRef.Get(ctx, &rawResult)
	if err != nil {
		return nil, err
	}

	for _, v := range rawResult {
		products = append(products, &v)
	}

	if products == nil {
		return nil, errors.New("No items")
	}
	return products, nil
}

// Get returns product given its ID
func (f *firebaseRepo) Get(ctx context.Context, productID string) (*model.Detail, error) {
	client, err := f.fb.Database(ctx)
	if err != nil {
		return nil, err // internal server error
	}

	product := model.Detail{}
	if err := client.NewRef(fmt.Sprintf("details/%s", productID)).Get(ctx, &product); err != nil {
		return nil, err
	}
	return &product, nil
}
