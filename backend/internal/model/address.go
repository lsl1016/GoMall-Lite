package model

import "time"

type Address struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index;not null"`
	Receiver  string `gorm:"type:varchar(64);not null"`
	Phone     string `gorm:"type:varchar(32);not null"`
	Province  string `gorm:"type:varchar(64);not null"`
	City      string `gorm:"type:varchar(64);not null"`
	District  string `gorm:"type:varchar(64);not null"`
	Detail    string `gorm:"type:varchar(255);not null"`
	IsDefault bool   `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ListAddresses(userID uint) ([]Address, error) {
	var addresses []Address
	if err := DB.Where("user_id = ?", userID).Order("is_default desc, id desc").Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func FindAddressByID(userID uint, id uint) (*Address, error) {
	var address Address
	if err := DB.Where("user_id = ? AND id = ?", userID, id).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func CreateAddress(address *Address) error {
	return DB.Create(address).Error
}

func SaveAddress(address *Address) error {
	return DB.Save(address).Error
}

func DeleteAddress(userID uint, id uint) error {
	return DB.Where("user_id = ? AND id = ?", userID, id).Delete(&Address{}).Error
}
