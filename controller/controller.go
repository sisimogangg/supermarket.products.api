package controller

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sisimogangg/supermarket.products.api/product"
	u "github.com/sisimogangg/supermarket.products.api/utils"
)

type productHandler struct {
	ProductService product.Service
}

//NewProductHandler creates a new instance of the product controller
func NewProductHandler(router *mux.Router, service product.Service) {
	handler := &productHandler{
		ProductService: service,
	}

	router.HandleFunc("/api/products", handler.allProducts).Methods("GET")

}

func (h *productHandler) allProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx != nil {
		ctx = context.Background()
	}

	products, err := h.ProductService.AllProducts(ctx)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
	}

	resp := u.Message(true, "success")
	resp["products"] = products

	u.Respond(w, resp)

}

func (h *productHandler) productByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx != nil {
		ctx = context.Background()
	}

	product, err := h.productByID(ctx,);
}
