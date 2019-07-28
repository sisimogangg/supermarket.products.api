package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/sisimogangg/supermarket.products.api/model"
	"github.com/sisimogangg/supermarket.products.api/repository"
	"github.com/sisimogangg/supermarket.products.api/utils"
)

type productService struct {
	repo    repository.DataAccessLayer
	timeout time.Duration
}

type discountCheck struct {
	IsOnDiscount bool   `json:"isondiscount"`
	ProductID    int32  `json:"productId"`
	Message      string `json:"message"`
	Status       bool   `json:"status"`
}

// NewProductService constructs a product service instance
func NewProductService(repo repository.DataAccessLayer, timeout time.Duration) Service {
	return &productService{repo, timeout}
}

func checkDiscountsOnServer(ctx context.Context, products []*model.Product) <-chan io.Reader {
	chanReaders := make(chan io.Reader)

	var wg sync.WaitGroup
	for _, p := range products {
		p := p // avoid capturing this
		wg.Add(1)
		go func() {
			URL := fmt.Sprintf("%s%v", viper.GetString("discount.verify"), p.ID)

			body, err := utils.GetRequest(ctx, URL)
			if err != nil {
				log.Fatal(err)
				return
			}

			select {
			case chanReaders <- *body:
			case <-ctx.Done():
				log.Fatal(errors.New("request Cancelled"))
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(chanReaders)
	}()
	return chanReaders
}

func readResponse(ctx context.Context, chanReaders <-chan io.Reader, chanDis chan<- discountCheck) {
	for r := range chanReaders {
		results := discountCheck{}

		err := json.NewDecoder(r).Decode(&results)
		if err != nil {
			log.Fatal(err)
		}

		select {
		case chanDis <- results:
		case <-ctx.Done():
			log.Fatal(errors.New("Request Cancelled"))
			return
		}
	}
}

func (s *productService) checkForProductDiscounts(ctx context.Context, products []*model.Product) ([]*model.Product, error) {

	//mapResponseToProductIDs := map[int32]interface{}{}
	mapResponseToProductIDs := make(map[int32]bool)
	for _, p := range products {
		mapResponseToProductIDs[p.ID] = false
	}

	chanReaders := checkDiscountsOnServer(ctx, products)
	chanDis := make(chan discountCheck)

	var wg sync.WaitGroup
	const num = 10

	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			readResponse(ctx, chanReaders, chanDis)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(chanDis)
	}()

	for d := range chanDis {
		mapResponseToProductIDs[d.ProductID] = d.IsOnDiscount
	}

	for _, p := range products {
		if isOnDiscount, ok := mapResponseToProductIDs[p.ID]; ok {
			p.Discount = bool(isOnDiscount)
		}
	}
	return products, nil
}

func (s *productService) List(ctx context.Context) ([]*model.Product, error) {
	time.Sleep(5 * time.Second)

	ps, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	/*	products, err := s.checkForProductDiscounts(ctx, ps)
		if err != nil {
			return nil, err
		}

		return products, nil */
	return ps, nil
}

func (s *productService) Get(ctx context.Context, productID string) (*model.Detail, error) {
	p, err := s.repo.Get(ctx, productID)
	if err != nil {
		return nil, err
	}

	return p, nil
}
