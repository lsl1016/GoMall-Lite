package service

import (
	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/model"
)

type ProductService struct{}

func NewProductService() *ProductService { return &ProductService{} }

func (s *ProductService) List(category string, keyword string) ([]dto.ProductDTO, error) {
	products, err := model.ListProducts(category, keyword)
	if err != nil {
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
		return nil, NewError(404, "商品不存在")
	}
	result := productDTO(p)
	return &result, nil
}

func productDTO(p *model.Product) dto.ProductDTO {
	return dto.ProductDTO{ID: p.ID, Name: p.Name, Price: p.Price, Stock: p.Stock, Image: p.Image, Category: p.Category, Description: p.Description}
}
