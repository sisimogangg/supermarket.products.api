package utils

import "github.com/sisimogangg/supermarket.products.api/model"

// Products is a seed of products
var Products = []*model.Product{
	&model.Product{
		ID:       100,
		ImageURL: "https://images.unsplash.com/photo-1478004521390-655bd10c9f43?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
		Name:     "Apple",
		Price: struct {
			Symbol   string `json:"symbol"`
			Currency string `json:"currency"`
			Amount   string `json:"amount"`
		}{
			Symbol:   "R",
			Currency: "RSA",
			Amount:   "2.00",
		},
	},
	&model.Product{
		ID:       200,
		ImageURL: "https://images.unsplash.com/photo-1528825871115-3581a5387919?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=658&q=80",
		Name:     "Banana",
		Price: struct {
			Symbol   string `json:"symbol"`
			Currency string `json:"currency"`
			Amount   string `json:"amount"`
		}{
			Symbol:   "R",
			Currency: "RSA",
			Amount:   "3.00",
		},
	},
	&model.Product{
		ID:       300,
		ImageURL: "https://images.unsplash.com/photo-1560769680-ba2f3767c785?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjExMDk0fQ&auto=format&fit=crop&w=700&q=80",
		Name:     "Coconut",
		Price: struct {
			Symbol   string `json:"symbol"`
			Currency string `json:"currency"`
			Amount   string `json:"amount"`
		}{
			Symbol:   "R",
			Currency: "RSA",
			Amount:   "4.00",
		},
	},
}

// Details stores product details seeding data
var Details = []*model.Detail{
	&model.Detail{
		Description: "Class 1, Top Red Apples, sweet and crispy.",
		Product:     *Products[0],
	},
	&model.Detail{
		Description: "Class 1, Easy-to-peel fresh, ripe Bananas.",
		Product:     *Products[1],
	},
	&model.Detail{
		Description: "Our coconuts can be eaten fresh or used for baking",
		Product:     *Products[2],
	},
}
