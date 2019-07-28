package model

// Product defines product item
type Product struct {
	ID       int32       `json:"id"`
	Name     string      `json:"name"`
	Discount interface{} `json:"discount"`
	ImageURL string      `json:"imageUrl"`
	Price    struct {
		Symbol   string `json:"symbol"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"price"`
}

// Detail stores product detials
type Detail struct {
	Description string `json:"description"`
	Discount    struct {
		DiscountID string `json:"discountID"`
		Summary    string `json:"summary"`
	}
	Product
}
