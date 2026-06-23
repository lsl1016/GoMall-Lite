package model

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gomall-lite-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg config.Config) error {
	var err error
	for i := 1; i <= 30; i++ {
		DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("waiting for mysql (%d/30): %v", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return err
	}

	if err := DB.AutoMigrate(&User{}, &Product{}, &CartItem{}, &Address{}, &Order{}, &OrderItem{}); err != nil {
		return err
	}

	return SeedData()
}

func SeedData() error {
	var user User
	if err := DB.Where("username = ?", "admin").First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		user = User{Username: "admin", PasswordHash: string(hash), Nickname: "张三"}
		if err := DB.Create(&user).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var count int64
	DB.Model(&Product{}).Count(&count)
	if count == 0 {
		products := []Product{
			{Name: "iPhone 15 Pro 256GB", Price: 7999, Stock: 100, Image: "/images/iphone.svg", Category: "手机数码", Description: "搭载 A17 Pro 芯片，钛金属边框，6.1 英寸超视网膜 XDR 显示屏，适合日常学习、办公和娱乐。"},
			{Name: "MacBook Air M2", Price: 8999, Stock: 50, Image: "/images/macbook.svg", Category: "电脑办公", Description: "轻薄便携的 M2 芯片笔记本，适合开发、办公、学习和内容创作。"},
			{Name: "索尼 WH-1000XM5", Price: 2599, Stock: 80, Image: "/images/headphone.svg", Category: "手机数码", Description: "旗舰级主动降噪耳机，拥有舒适佩戴体验和优秀音质表现。"},
			{Name: "小米 14", Price: 4299, Stock: 120, Image: "/images/mi14.svg", Category: "手机数码", Description: "高性能安卓旗舰手机，轻薄机身，日常拍照和游戏体验优秀。"},
			{Name: "SK-II 神仙水 230ml", Price: 1230, Stock: 90, Image: "/images/sk2.svg", Category: "美妆个护", Description: "经典护肤精华水，适合日常护肤使用。"},
			{Name: "Nike Air Force 1", Price: 699, Stock: 60, Image: "/images/shoe.svg", Category: "服饰鞋包", Description: "经典百搭休闲鞋，适合多场景穿搭。"},
			{Name: "Apple Watch Series 9", Price: 2999, Stock: 70, Image: "/images/watch.svg", Category: "手机数码", Description: "健康记录、运动追踪和消息提醒的智能手表。"},
			{Name: "AirPods Pro 第二代", Price: 1899, Stock: 85, Image: "/images/airpods.svg", Category: "手机数码", Description: "主动降噪真无线耳机，便携易用。"},
			{Name: "三只松鼠零食礼包", Price: 129, Stock: 200, Image: "/images/snack.svg", Category: "食品生鲜", Description: "多口味零食组合，适合家庭休闲和办公室分享。"},
			{Name: "美的空气炸锅", Price: 399, Stock: 45, Image: "/images/airfryer.svg", Category: "家用电器", Description: "便捷厨房小家电，轻松制作低油美食。"},
		}
		if err := DB.Create(&products).Error; err != nil {
			return err
		}
	}

	if user.ID != 0 {
		var addrCount int64
		DB.Model(&Address{}).Where("user_id = ?", user.ID).Count(&addrCount)
		if addrCount == 0 {
			addr := Address{UserID: user.ID, Receiver: "张三", Phone: "13800138000", Province: "北京市", City: "北京市", District: "朝阳区", Detail: "望京街道 88 号 6 号楼", IsDefault: true}
			if err := DB.Create(&addr).Error; err != nil {
				return err
			}
		}
	}

	fmt.Println("seed data ready")
	return nil
}
