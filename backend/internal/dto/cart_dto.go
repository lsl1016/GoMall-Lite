package dto

type AddCartRequest struct {
	ProductID uint `json:"productId" binding:"required"`
	Count     int  `json:"count"`
}

type UpdateCartRequest struct {
	Count   *int  `json:"count"`
	Checked *bool `json:"checked"`
}

type CartItemDTO struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"productId"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Stock     int    `json:"stock"`
	Image     string `json:"image"`
	Count     int    `json:"count"`
	Checked   bool   `json:"checked"`
	Subtotal  int    `json:"subtotal"`
}
