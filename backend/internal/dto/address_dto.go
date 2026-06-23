package dto

type AddressRequest struct {
	Receiver  string `json:"receiver" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Province  string `json:"province" binding:"required"`
	City      string `json:"city" binding:"required"`
	District  string `json:"district" binding:"required"`
	Detail    string `json:"detail" binding:"required"`
	IsDefault bool   `json:"isDefault"`
}

type AddressDTO struct {
	ID        uint   `json:"id"`
	Receiver  string `json:"receiver"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Detail    string `json:"detail"`
	IsDefault bool   `json:"isDefault"`
}
