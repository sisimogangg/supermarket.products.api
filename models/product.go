package models

// Product defines product item
type Product struct {
	ID           int32
	Name         string
	IsOnDiscount bool
	ImageURL     string
	Price        struct {
		Symbol   rune
		Currency string
		Amount   float32
	}
}
