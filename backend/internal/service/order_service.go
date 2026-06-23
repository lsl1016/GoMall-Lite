package service

import (
	"fmt"
	"time"

	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/logger"
	"gomall-lite-api/internal/model"
	"gorm.io/gorm"
)

const (
	OrderPending   = "待支付"
	OrderPaid      = "已支付"
	OrderCancelled = "已取消"
	OrderCompleted = "已完成"
)

type OrderService struct{}

func NewOrderService() *OrderService { return &OrderService{} }

type orderCreateItem struct {
	CartID  uint
	Product model.Product
	Count   int
}

func (s *OrderService) Create(userID uint, req dto.CreateOrderRequest) (*dto.OrderDTO, error) {
	logger.Default().Info("create order start", "user_id", userID, "address_id", req.AddressID, "item_count", len(req.Items))
	address, err := model.FindAddressByID(userID, req.AddressID)
	if err != nil {
		logger.Default().Warn("create order failed: address not found", "user_id", userID, "address_id", req.AddressID)
		return nil, NewError(404, "收货地址不存在")
	}

	items, err := s.buildOrderItems(userID, req.Items)
	if err != nil {
		logger.Default().Warn("create order failed: build items failed", "user_id", userID, "error", err)
		return nil, err
	}
	if len(items) == 0 {
		logger.Default().Warn("create order failed: empty items", "user_id", userID)
		return nil, NewError(400, "请选择要结算的商品")
	}

	total := 0
	for _, item := range items {
		if item.Count <= 0 {
			logger.Default().Warn("create order failed: invalid item count", "user_id", userID, "product_id", item.Product.ID, "count", item.Count)
			return nil, NewError(400, "商品数量错误")
		}
		if item.Product.Stock < item.Count {
			logger.Default().Warn("create order failed: stock not enough", "user_id", userID, "product_id", item.Product.ID, "stock", item.Product.Stock, "count", item.Count)
			return nil, NewError(400, fmt.Sprintf("%s 库存不足", item.Product.Name))
		}
		total += item.Product.Price * item.Count
	}

	addressSnapshot := fmt.Sprintf("%s %s，%s%s%s%s", address.Receiver, address.Phone, address.Province, address.City, address.District, address.Detail)
	orderNo := fmt.Sprintf("ORDER%s%04d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%10000)
	order := model.Order{UserID: userID, OrderNo: orderNo, TotalAmount: total, Status: OrderPending, AddressSnapshot: addressSnapshot, Remark: req.Remark}

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		for _, item := range items {
			res := tx.Model(&model.Product{}).Where("id = ? AND stock >= ?", item.Product.ID, item.Count).Update("stock", gorm.Expr("stock - ?", item.Count))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return NewError(400, fmt.Sprintf("%s 库存不足", item.Product.Name))
			}

			orderItem := model.OrderItem{OrderID: order.ID, ProductID: item.Product.ID, ProductName: item.Product.Name, ProductImage: item.Product.Image, Price: item.Product.Price, Count: item.Count, Subtotal: item.Product.Price * item.Count}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}
			if item.CartID != 0 {
				if err := tx.Where("user_id = ? AND id = ?", userID, item.CartID).Delete(&model.CartItem{}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		logger.Default().Error("create order transaction failed", "user_id", userID, "order_no", orderNo, "error", err)
		return nil, err
	}

	created, err := model.FindOrderByID(userID, order.ID)
	if err != nil {
		logger.Default().Error("create order reload failed", "user_id", userID, "order_id", order.ID, "error", err)
		return nil, err
	}
	result := orderDTO(created)
	logger.Default().Info("create order success", "user_id", userID, "order_id", order.ID, "order_no", order.OrderNo, "total_amount", order.TotalAmount)
	return &result, nil
}

func (s *OrderService) buildOrderItems(userID uint, reqItems []dto.CreateOrderItemRequest) ([]orderCreateItem, error) {
	items := make([]orderCreateItem, 0)
	if len(reqItems) == 0 {
		var cartItems []model.CartItem
		if err := model.DB.Preload("Product").Where("user_id = ? AND checked = ?", userID, true).Find(&cartItems).Error; err != nil {
			logger.Default().Error("build order items from cart failed", "user_id", userID, "error", err)
			return nil, err
		}
		for _, cart := range cartItems {
			items = append(items, orderCreateItem{CartID: cart.ID, Product: cart.Product, Count: cart.Count})
		}
		return items, nil
	}

	for _, req := range reqItems {
		if req.CartID != 0 {
			cart, err := model.FindCartItemByID(userID, req.CartID)
			if err != nil {
				logger.Default().Warn("build order items failed: cart item not found", "user_id", userID, "cart_item_id", req.CartID)
				return nil, NewError(404, "购物车商品不存在")
			}
			count := cart.Count
			if req.Count > 0 {
				count = req.Count
			}
			items = append(items, orderCreateItem{CartID: cart.ID, Product: cart.Product, Count: count})
			continue
		}
		product, err := model.GetProductByID(req.ProductID)
		if err != nil {
			logger.Default().Warn("build order items failed: product not found", "user_id", userID, "product_id", req.ProductID)
			return nil, NewError(404, "商品不存在")
		}
		items = append(items, orderCreateItem{Product: *product, Count: req.Count})
	}
	return items, nil
}

func (s *OrderService) List(userID uint) ([]dto.OrderDTO, error) {
	orders, err := model.ListOrders(userID)
	if err != nil {
		logger.Default().Error("list orders failed", "user_id", userID, "error", err)
		return nil, err
	}
	result := make([]dto.OrderDTO, 0, len(orders))
	for _, order := range orders {
		result = append(result, orderDTO(&order))
	}
	return result, nil
}

func (s *OrderService) Detail(userID uint, id uint) (*dto.OrderDTO, error) {
	order, err := model.FindOrderByID(userID, id)
	if err != nil {
		logger.Default().Warn("order detail failed: order not found", "user_id", userID, "order_id", id)
		return nil, NewError(404, "订单不存在")
	}
	result := orderDTO(order)
	return &result, nil
}

func (s *OrderService) Pay(userID uint, id uint) (*dto.OrderDTO, error) {
	order, err := model.FindOrderByID(userID, id)
	if err != nil {
		logger.Default().Warn("pay order failed: order not found", "user_id", userID, "order_id", id)
		return nil, NewError(404, "订单不存在")
	}
	if order.Status != OrderPending {
		logger.Default().Warn("pay order failed: invalid status", "user_id", userID, "order_id", id, "status", order.Status)
		return nil, NewError(400, "只有待支付订单可以支付")
	}
	order.Status = OrderPaid
	if err := model.SaveOrder(order); err != nil {
		logger.Default().Error("pay order save failed", "user_id", userID, "order_id", id, "error", err)
		return nil, err
	}
	result := orderDTO(order)
	logger.Default().Info("pay order success", "user_id", userID, "order_id", id, "order_no", order.OrderNo)
	return &result, nil
}

func (s *OrderService) Cancel(userID uint, id uint) (*dto.OrderDTO, error) {
	order, err := model.FindOrderByID(userID, id)
	if err != nil {
		logger.Default().Warn("cancel order failed: order not found", "user_id", userID, "order_id", id)
		return nil, NewError(404, "订单不存在")
	}
	if order.Status != OrderPending {
		logger.Default().Warn("cancel order failed: invalid status", "user_id", userID, "order_id", id, "status", order.Status)
		return nil, NewError(400, "只有待支付订单可以取消")
	}
	order.Status = OrderCancelled
	if err := model.SaveOrder(order); err != nil {
		logger.Default().Error("cancel order save failed", "user_id", userID, "order_id", id, "error", err)
		return nil, err
	}
	result := orderDTO(order)
	logger.Default().Info("cancel order success", "user_id", userID, "order_id", id, "order_no", order.OrderNo)
	return &result, nil
}

func orderDTO(order *model.Order) dto.OrderDTO {
	items := make([]dto.OrderItemDTO, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, dto.OrderItemDTO{ID: item.ID, ProductID: item.ProductID, ProductName: item.ProductName, ProductImage: item.ProductImage, Price: item.Price, Count: item.Count, Subtotal: item.Subtotal})
	}
	return dto.OrderDTO{ID: order.ID, OrderNo: order.OrderNo, TotalAmount: order.TotalAmount, Status: order.Status, AddressSnapshot: order.AddressSnapshot, Remark: order.Remark, CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"), Items: items}
}
