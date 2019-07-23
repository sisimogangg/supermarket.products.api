package models

// Product defines product item
type Product struct {
	ID        int
	Name      string
	Promotion string
	ImageURL  string
	Price     struct {
		Symbol   rune
		Currency string
		Amount   float32
	}
}
