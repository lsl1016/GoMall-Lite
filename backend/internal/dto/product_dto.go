package dto

type ProductDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Image       string `json:"image"`
	Category    string `json:"category"`
	Description string `json:"description"`
}
