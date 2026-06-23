package model

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"type:varchar(64);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Nickname     string `gorm:"type:varchar(64)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CreateUser(user *User) error {
	return DB.Create(user).Error
}

func FindUserByUsername(username string) (*User, error) {
	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByID(id uint) (*User, error) {
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
