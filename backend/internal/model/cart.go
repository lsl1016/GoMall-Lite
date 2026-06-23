package model

import "time"

type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index;not null"`
	ProductID uint `gorm:"index;not null"`
	Product   Product
	Count     int  `gorm:"not null"`
	Checked   bool `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ListCartItems(userID uint) ([]CartItem, error) {
	var items []CartItem
	if err := DB.Preload("Product").Where("user_id = ?", userID).Order("id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func FindCartItemByUserProduct(userID uint, productID uint) (*CartItem, error) {
	var item CartItem
	if err := DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func FindCartItemByID(userID uint, id uint) (*CartItem, error) {
	var item CartItem
	if err := DB.Preload("Product").Where("user_id = ? AND id = ?", userID, id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateCartItem(item *CartItem) error {
	return DB.Create(item).Error
}

func SaveCartItem(item *CartItem) error {
	return DB.Save(item).Error
}

func DeleteCartItem(userID uint, id uint) error {
	return DB.Where("user_id = ? AND id = ?", userID, id).Delete(&CartItem{}).Error
}

func ClearCart(userID uint) error {
	return DB.Where("user_id = ?", userID).Delete(&CartItem{}).Error
}
