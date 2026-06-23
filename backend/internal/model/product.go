package model

import "time"

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(128);not null"`
	Price       int    `gorm:"not null"`
	Stock       int    `gorm:"not null"`
	Image       string `gorm:"type:varchar(255)"`
	Category    string `gorm:"type:varchar(64);index"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ListProducts(category string, keyword string) ([]Product, error) {
	var products []Product
	query := DB.Model(&Product{})
	if category != "" && category != "全部" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if err := query.Order("id asc").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func GetProductByID(id uint) (*Product, error) {
	var product Product
	if err := DB.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
