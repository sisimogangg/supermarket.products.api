package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/sisimogangg/supermarket.products.api/models"
	"github.com/sisimogangg/supermarket.products.api/repository"
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

func contactDiscountServer(ctx context.Context, products []*models.Product) <-chan io.Reader {
	chanReaders := make(chan io.Reader)
	//defer close(chanReaders)
	fmt.Println("contacting discount server")

	var wg sync.WaitGroup
	for _, p := range products {
		p := p // avoid capturing this
		wg.Add(1)
		go func() {
			fmt.Println("fetching discount for : ", p.ID)
			URL := fmt.Sprintf("%s%v", viper.GetString("discount.verify"), p.ID)

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			if err != nil {
				log.Fatal(err)
				return
			}

			req = req.WithContext(ctx)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
				return
			}

			if res.StatusCode != http.StatusOK {
				log.Fatal(errors.New(res.Status))
			}

			select {
			case chanReaders <- res.Body:
				fmt.Println("fetched for : ", p.ID)
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

func (s *productService) retrieveDiscounts(ctx context.Context, products []*models.Product) ([]*models.Product, error) {

	//mapResponseToProductIDs := map[int32]interface{}{}
	mapResponseToProductIDs := make(map[int32]bool)
	for _, p := range products {
		mapResponseToProductIDs[p.ID] = false
	}

	chanReaders := contactDiscountServer(ctx, products)
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
			p.IsOnDiscount = isOnDiscount
		}
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
