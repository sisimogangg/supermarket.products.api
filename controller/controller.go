package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sisimogangg/supermarket.products.api/service"
	u "github.com/sisimogangg/supermarket.products.api/utils"
)

type productHandler struct {
	productService service.Service
}

//NewProductHandler creates a new instance of the product controller
func NewProductHandler(router *mux.Router, service service.Service) {
	handler := &productHandler{
		productService: service,
	}

	router.HandleFunc("/api/products", handler.allProducts).Methods("GET")
	router.HandleFunc("/api/products/{id}", handler.getProductByID).Methods("GET")

}

func (h *productHandler) allProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx != nil {
		ctx = context.Background()
	}

	products, err := h.productService.AllProducts(ctx)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	resp := u.Message(true, "success")
	resp["products"] = products

	u.Respond(w, resp)

}

func (h *productHandler) getProductByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	ctx := r.Context()
	if ctx != nil {
		ctx = context.Background()
	}

	product, err := h.productService.GetProductByID(ctx, int32(productID))
	if err != nil {
		if errVal, ok := err.(*u.HTTPError); ok {
			w.WriteHeader(errVal.Status)
		}
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	resp := u.Message(true, "success")
	resp["product"] = product

	u.Respond(w, resp)

}
