package repository

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go"
	pb "github.com/sisimogangg/supermarket.products.api/proto"
)

// DataAccessLayer defines repository expected behaviour
type DataAccessLayer interface {
	List(ctx context.Context) ([]*pb.Product, error)
	Get(ctx context.Context, productID string) (*pb.ProductDetail, error)
}

type firebaseRepo struct {
	fb *firebase.App
}

// NewFirebaseRepo defines a constructor for firebaserepo
func NewFirebaseRepo(app *firebase.App) DataAccessLayer {
	return &firebaseRepo{app}
}

//List returns all products
func (f *firebaseRepo) List(ctx context.Context) ([]*pb.Product, error) {
	client, err := f.fb.Database(ctx)
	productsRef := client.NewRef("products")
	if err != nil {
		return nil, err // internal server error
	}

	products := []*pb.Product{}
	var rawResult map[string]pb.Product

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
func (f *firebaseRepo) Get(ctx context.Context, productID string) (*pb.ProductDetail, error) {
	client, err := f.fb.Database(ctx)
	if err != nil {
		return nil, err // internal server error
	}

	product := pb.ProductDetail{}
	if err := client.NewRef(fmt.Sprintf("details/%s", productID)).Get(ctx, &product); err != nil {
		return nil, err
	}
	return &product, nil
}
