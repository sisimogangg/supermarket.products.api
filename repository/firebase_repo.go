package repository

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	pb "github.com/sisimogangg/supermarket.products.api/proto"
)

// Repository defines repository expected behaviour
type Repository interface {
	List(ctx context.Context) ([]*pb.Product, error)
	Get(ctx context.Context, productID string) (*pb.ProductDetail, error)
}

type firebaseRepo struct {
	fb *firebase.App
}

// NewFirebaseRepo defines a constructor for firebaserepo
func NewFirebaseRepo(app *firebase.App) Repository {
	return &firebaseRepo{app}
}

//List returns all products
func (f *firebaseRepo) List(ctx context.Context) ([]*pb.Product, error) {
	client, err := f.fb.Database(ctx)
	if err != nil {
		return nil, err // internal server error
	}
	productsRef := client.NewRef("products")

	products := []*pb.Product{}
	var rawResult map[string]pb.Product

	err = productsRef.Get(ctx, &rawResult)
	if err != nil {
		return nil, err
	}

	for _, v := range rawResult {
		products = append(products, &v)
	}

	return products, nil
}

// Get returns product given its ID
func (f *firebaseRepo) Get(ctx context.Context, productID string) (*pb.ProductDetail, error) {
	client, err := f.fb.Database(ctx)
	if err != nil {
		return nil, err // internal server error
	}
	pDetailsRef := client.NewRef(fmt.Sprintf("details/%s", productID))

	pDetail := pb.ProductDetail{}
	if err := pDetailsRef.Get(ctx, &pDetail); err != nil {
		return nil, err
	}
	return &pDetail, nil
}
