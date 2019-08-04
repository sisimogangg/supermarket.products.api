package utils

import pb "github.com/sisimogangg/supermarket.products.api/proto"

var price1 = pb.Price{
	Symbol:   "R",
	Currency: "RSA",
	Amount:   "2.00",
}

var price2 = pb.Price{
	Symbol:   "R",
	Currency: "RSA",
	Amount:   "3.00",
}

var price3 = pb.Price{
	Symbol:   "R",
	Currency: "RSA",
	Amount:   "4.00",
}

// Products is a seed of products
var Products = []*pb.Product{
	&pb.Product{
		Id:       "100",
		ImageURL: "https://images.unsplash.com/photo-1478004521390-655bd10c9f43?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
		Name:     "Apple",
		Discount: false,
		Price:    &price1,
	},
	&pb.Product{
		Id:       "200",
		ImageURL: "https://images.unsplash.com/photo-1528825871115-3581a5387919?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=658&q=80",
		Name:     "Banana",
		Discount: false,
		Price:    &price2,
	},
	&pb.Product{
		Id:       "300",
		ImageURL: "https://images.unsplash.com/photo-1560769680-ba2f3767c785?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjExMDk0fQ&auto=format&fit=crop&w=700&q=80",
		Name:     "Coconut",
		Discount: false,
		Price:    &price3,
	},
}

// Details stores product details seeding data
var Details = []*pb.ProductDetail{
	&pb.ProductDetail{
		Description: "Class 1, Top Red Apples, sweet and crispy.",
		Discount:    nil,
		Product:     Products[0],
	},
	&pb.ProductDetail{
		Description: "Class 1, Easy-to-peel fresh, ripe Bananas.",
		Product:     Products[1],
	},
	&pb.ProductDetail{
		Description: "Our coconuts can be eaten fresh or used for baking",
		Product:     Products[2],
	},
}
