package dto

type CreateOrderItemRequest struct {
	CartID    uint `json:"cartId"`
	ProductID uint `json:"productId"`
	Count     int  `json:"count"`
}

type CreateOrderRequest struct {
	AddressID uint                     `json:"addressId" binding:"required"`
	Remark    string                   `json:"remark"`
	Items     []CreateOrderItemRequest `json:"items"`
}

type OrderItemDTO struct {
	ID           uint   `json:"id"`
	ProductID    uint   `json:"productId"`
	ProductName  string `json:"productName"`
	ProductImage string `json:"productImage"`
	Price        int    `json:"price"`
	Count        int    `json:"count"`
	Subtotal     int    `json:"subtotal"`
}

type OrderDTO struct {
	ID              uint           `json:"id"`
	OrderNo         string         `json:"orderNo"`
	TotalAmount     int            `json:"totalAmount"`
	Status          string         `json:"status"`
	AddressSnapshot string         `json:"addressSnapshot"`
	Remark          string         `json:"remark"`
	CreatedAt       string         `json:"createdAt"`
	Items           []OrderItemDTO `json:"items"`
}
