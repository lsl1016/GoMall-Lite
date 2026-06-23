package model

import "time"

type Order struct {
	ID              uint   `gorm:"primaryKey"`
	UserID          uint   `gorm:"index;not null"`
	OrderNo         string `gorm:"type:varchar(64);uniqueIndex;not null"`
	TotalAmount     int    `gorm:"not null"`
	Status          string `gorm:"type:varchar(32);not null"`
	AddressSnapshot string `gorm:"type:text"`
	Remark          string `gorm:"type:varchar(255)"`
	Items           []OrderItem
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ID           uint `gorm:"primaryKey"`
	OrderID      uint `gorm:"index;not null"`
	ProductID    uint
	ProductName  string `gorm:"type:varchar(128);not null"`
	ProductImage string `gorm:"type:varchar(255)"`
	Price        int    `gorm:"not null"`
	Count        int    `gorm:"not null"`
	Subtotal     int    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func ListOrders(userID uint) ([]Order, error) {
	var orders []Order
	if err := DB.Preload("Items").Where("user_id = ?", userID).Order("id desc").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func FindOrderByID(userID uint, id uint) (*Order, error) {
	var order Order
	if err := DB.Preload("Items").Where("user_id = ? AND id = ?", userID, id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func SaveOrder(order *Order) error {
	return DB.Save(order).Error
}
