package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"

	discountProto "github.com/sisimogangg/supermarket.discount.api/proto"
	pb "github.com/sisimogangg/supermarket.products.api/proto"

	"github.com/sisimogangg/supermarket.products.api/repository"
	"github.com/sisimogangg/supermarket.products.api/utils"
)

type productService struct {
	Repo           repository.Repository
	DiscountClient discountProto.DiscountServiceClient
	Timeout        time.Duration
}

// NewProductService creates and returns a new instance of productService
func NewProductService(repo repository.Repository, ds discountProto.DiscountServiceClient, timeout time.Duration) pb.ProductServiceHandler {
	return &productService{repo, ds, timeout}
}

type discountCheck struct {
	IsOnDiscount bool   `json:"isondiscount"`
	ProductID    int32  `json:"productId"`
	Message      string `json:"message"`
	Status       bool   `json:"status"`
}

func retrieveDiscountsFromServer(ctx context.Context, products []*pb.Product) <-chan io.Reader {
	chanReaders := make(chan io.Reader)

	var wg sync.WaitGroup
	for _, p := range products {
		p := p // avoid capturing this
		wg.Add(1)
		go func() {
			URL := fmt.Sprintf("%s%v", viper.GetString("discount.verify"), p.Id)

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

func (s *productService) getProductDiscounts(ctx context.Context, products []*pb.Product) ([]*pb.Product, error) {

	//mapResponseToProductIDs := map[int32]interface{}{}
	mapResponseToProductIDs := make(map[int32]bool)
	for _, p := range products {
		index, err := strconv.Atoi(p.Id)
		if err != nil {
			return nil, err
		}
		mapResponseToProductIDs[int32(index)] = false
	}

	chanReaders := retrieveDiscountsFromServer(ctx, products)
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
		index, err := strconv.Atoi(p.Id)
		if err != nil {
			return nil, err
		}
		if isOnDiscount, ok := mapResponseToProductIDs[int32(index)]; ok {
			p.Discount = bool(isOnDiscount)
		}
	}
	return products, nil
}

// List returns a list of products
func (s *productService) List(ctx context.Context, req *pb.ListRequest, resp *pb.ListResponse) error {

	ps, err := s.Repo.List(ctx)
	if err != nil {
		return err
	}

	// check if discount is available
	wg := sync.WaitGroup{}
	for _, p := range ps {
		p := p
		wg.Add(1)
		go func() {
			getReq := discountProto.GetRequest{ProductID: p.Id}
			discount, err := s.DiscountClient.Get(ctx, &getReq)
			if err != nil {
				log.Fatal(err)
			}

			if len(discount.DiscountID) > 0 {
				p.Discount = true
			} else {
				p.Discount = false
			}
			wg.Done()
		}()
	}
	wg.Wait()

	resp.Products = ps

	return nil
}

// Get returns product details
func (s *productService) Get(ctx context.Context, req *pb.GetRequest, resp *pb.ProductDetail) error {
	p, err := s.Repo.Get(ctx, req.Id)
	if err != nil {
		return err
	}

	resp.Description = p.Description

	getReq := discountProto.GetRequest{ProductID: req.Id}
	discount, err := s.DiscountClient.Get(ctx, &getReq)
	if err != nil {
		return nil
	}

	if len(discount.DiscountID) > 0 {
		resp.Discount.DiscountID = discount.DiscountID
		resp.Discount.Summary = discount.Summary
	}
	resp.Product = p.Product

	return nil
}
