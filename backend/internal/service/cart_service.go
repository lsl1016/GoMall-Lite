package service

import (
	"errors"

	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/model"
	"gorm.io/gorm"
)

type CartService struct{}

func NewCartService() *CartService { return &CartService{} }

func (s *CartService) List(userID uint) ([]dto.CartItemDTO, error) {
	items, err := model.ListCartItems(userID)
	if err != nil {
		return nil, err
	}
	return cartItemDTOs(items), nil
}

func (s *CartService) Add(userID uint, req dto.AddCartRequest) ([]dto.CartItemDTO, error) {
	if req.Count <= 0 {
		req.Count = 1
	}
	product, err := model.GetProductByID(req.ProductID)
	if err != nil {
		return nil, NewError(404, "商品不存在")
	}
	if product.Stock < req.Count {
		return nil, NewError(400, "商品库存不足")
	}

	item, err := model.FindCartItemByUserProduct(userID, req.ProductID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newItem := model.CartItem{UserID: userID, ProductID: req.ProductID, Count: req.Count, Checked: true}
		if err := model.CreateCartItem(&newItem); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		if item.Count+req.Count > product.Stock {
			return nil, NewError(400, "商品库存不足")
		}
		item.Count += req.Count
		item.Checked = true
		if err := model.SaveCartItem(item); err != nil {
			return nil, err
		}
	}
	return s.List(userID)
}

func (s *CartService) Update(userID uint, id uint, req dto.UpdateCartRequest) ([]dto.CartItemDTO, error) {
	item, err := model.FindCartItemByID(userID, id)
	if err != nil {
		return nil, NewError(404, "购物车商品不存在")
	}
	if req.Count != nil {
		if *req.Count <= 0 {
			return nil, NewError(400, "数量必须大于 0")
		}
		if *req.Count > item.Product.Stock {
			return nil, NewError(400, "商品库存不足")
		}
		item.Count = *req.Count
	}
	if req.Checked != nil {
		item.Checked = *req.Checked
	}
	if err := model.SaveCartItem(item); err != nil {
		return nil, err
	}
	return s.List(userID)
}

func (s *CartService) Remove(userID uint, id uint) ([]dto.CartItemDTO, error) {
	if err := model.DeleteCartItem(userID, id); err != nil {
		return nil, err
	}
	return s.List(userID)
}

func (s *CartService) Clear(userID uint) error {
	return model.ClearCart(userID)
}

func cartItemDTOs(items []model.CartItem) []dto.CartItemDTO {
	result := make([]dto.CartItemDTO, 0, len(items))
	for _, item := range items {
		result = append(result, dto.CartItemDTO{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Price:     item.Product.Price,
			Stock:     item.Product.Stock,
			Image:     item.Product.Image,
			Count:     item.Count,
			Checked:   item.Checked,
			Subtotal:  item.Product.Price * item.Count,
		})
	}
	return result
}
