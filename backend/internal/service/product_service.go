package service

import (
	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/logger"
	"gomall-lite-api/internal/model"
)

type ProductService struct{}

func NewProductService() *ProductService { return &ProductService{} }

func (s *ProductService) List(category string, keyword string) ([]dto.ProductDTO, error) {
	logger.Default().Debug("list products", "category", category, "keyword", keyword)
	products, err := model.ListProducts(category, keyword)
	if err != nil {
		logger.Default().Error("list products failed", "category", category, "keyword", keyword, "error", err)
		return nil, err
	}
	result := make([]dto.ProductDTO, 0, len(products))
	for _, p := range products {
		result = append(result, productDTO(&p))
	}
	return result, nil
}

func (s *ProductService) Detail(id uint) (*dto.ProductDTO, error) {
	p, err := model.GetProductByID(id)
	if err != nil {
		logger.Default().Warn("product detail failed: product not found", "product_id", id)
		return nil, NewError(404, "商品不存在")
	}
	result := productDTO(p)
	return &result, nil
}

func productDTO(p *model.Product) dto.ProductDTO {
	return dto.ProductDTO{ID: p.ID, Name: p.Name, Price: p.Price, Stock: p.Stock, Image: p.Image, Category: p.Category, Description: p.Description}
}
