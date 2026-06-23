package service

import (
	"errors"

	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/logger"
	"gomall-lite-api/internal/model"
	"gorm.io/gorm"
)

type CartService struct{}

func NewCartService() *CartService { return &CartService{} }

func (s *CartService) List(userID uint) ([]dto.CartItemDTO, error) {
	items, err := model.ListCartItems(userID)
	if err != nil {
		logger.Default().Error("list cart failed", "user_id", userID, "error", err)
		return nil, err
	}
	return cartItemDTOs(items), nil
}

func (s *CartService) Add(userID uint, req dto.AddCartRequest) ([]dto.CartItemDTO, error) {
	logger.Default().Info("add to cart", "user_id", userID, "product_id", req.ProductID, "count", req.Count)
	if req.Count <= 0 {
		req.Count = 1
	}

	product, err := model.GetProductByID(req.ProductID)
	if err != nil {
		logger.Default().Warn("add to cart failed: product not found", "user_id", userID, "product_id", req.ProductID)
		return nil, NewError(404, "商品不存在")
	}
	if product.Stock < req.Count {
		logger.Default().Warn("add to cart failed: stock not enough", "user_id", userID, "product_id", req.ProductID, "stock", product.Stock, "count", req.Count)
		return nil, NewError(400, "商品库存不足")
	}

	item, err := model.FindCartItemByUserProduct(userID, req.ProductID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newItem := model.CartItem{UserID: userID, ProductID: req.ProductID, Count: req.Count, Checked: true}
		if err := model.CreateCartItem(&newItem); err != nil {
			logger.Default().Error("add to cart create failed", "user_id", userID, "product_id", req.ProductID, "error", err)
			return nil, err
		}
	} else if err != nil {
		logger.Default().Error("add to cart lookup failed", "user_id", userID, "product_id", req.ProductID, "error", err)
		return nil, err
	} else {
		if item.Count+req.Count > product.Stock {
			logger.Default().Warn("add to cart failed: total stock not enough", "user_id", userID, "product_id", req.ProductID, "stock", product.Stock, "count", item.Count+req.Count)
			return nil, NewError(400, "商品库存不足")
		}
		item.Count += req.Count
		item.Checked = true
		if err := model.SaveCartItem(item); err != nil {
			logger.Default().Error("add to cart save failed", "user_id", userID, "product_id", req.ProductID, "cart_item_id", item.ID, "error", err)
			return nil, err
		}
	}

	logger.Default().Info("add to cart success", "user_id", userID, "product_id", req.ProductID)
	return s.List(userID)
}

func (s *CartService) Update(userID uint, id uint, req dto.UpdateCartRequest) ([]dto.CartItemDTO, error) {
	item, err := model.FindCartItemByID(userID, id)
	if err != nil {
		logger.Default().Warn("update cart failed: item not found", "user_id", userID, "cart_item_id", id)
		return nil, NewError(404, "购物车商品不存在")
	}
	if req.Count != nil {
		if *req.Count <= 0 {
			logger.Default().Warn("update cart failed: invalid count", "user_id", userID, "cart_item_id", id, "count", *req.Count)
			return nil, NewError(400, "数量必须大于 0")
		}
		if *req.Count > item.Product.Stock {
			logger.Default().Warn("update cart failed: stock not enough", "user_id", userID, "cart_item_id", id, "stock", item.Product.Stock, "count", *req.Count)
			return nil, NewError(400, "商品库存不足")
		}
		item.Count = *req.Count
	}
	if req.Checked != nil {
		item.Checked = *req.Checked
	}
	if err := model.SaveCartItem(item); err != nil {
		logger.Default().Error("update cart save failed", "user_id", userID, "cart_item_id", id, "error", err)
		return nil, err
	}

	logger.Default().Info("update cart success", "user_id", userID, "cart_item_id", id)
	return s.List(userID)
}

func (s *CartService) Remove(userID uint, id uint) ([]dto.CartItemDTO, error) {
	if err := model.DeleteCartItem(userID, id); err != nil {
		logger.Default().Error("remove cart item failed", "user_id", userID, "cart_item_id", id, "error", err)
		return nil, err
	}
	logger.Default().Info("remove cart item success", "user_id", userID, "cart_item_id", id)
	return s.List(userID)
}

func (s *CartService) Clear(userID uint) error {
	if err := model.ClearCart(userID); err != nil {
		logger.Default().Error("clear cart failed", "user_id", userID, "error", err)
		return err
	}
	logger.Default().Info("clear cart success", "user_id", userID)
	return nil
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
