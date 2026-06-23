package service

import (
	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/model"
)

type AddressService struct{}

func NewAddressService() *AddressService { return &AddressService{} }

func (s *AddressService) List(userID uint) ([]dto.AddressDTO, error) {
	addresses, err := model.ListAddresses(userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.AddressDTO, 0, len(addresses))
	for _, addr := range addresses {
		result = append(result, addressDTO(&addr))
	}
	return result, nil
}

func (s *AddressService) Create(userID uint, req dto.AddressRequest) ([]dto.AddressDTO, error) {
	if req.IsDefault {
		if err := model.DB.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
			return nil, err
		}
	}
	addr := model.Address{UserID: userID, Receiver: req.Receiver, Phone: req.Phone, Province: req.Province, City: req.City, District: req.District, Detail: req.Detail, IsDefault: req.IsDefault}
	if err := model.CreateAddress(&addr); err != nil {
		return nil, err
	}
	return s.List(userID)
}

func (s *AddressService) Update(userID uint, id uint, req dto.AddressRequest) ([]dto.AddressDTO, error) {
	addr, err := model.FindAddressByID(userID, id)
	if err != nil {
		return nil, NewError(404, "地址不存在")
	}
	if req.IsDefault {
		if err := model.DB.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
			return nil, err
		}
	}
	addr.Receiver = req.Receiver
	addr.Phone = req.Phone
	addr.Province = req.Province
	addr.City = req.City
	addr.District = req.District
	addr.Detail = req.Detail
	addr.IsDefault = req.IsDefault
	if err := model.SaveAddress(addr); err != nil {
		return nil, err
	}
	return s.List(userID)
}

func (s *AddressService) Delete(userID uint, id uint) ([]dto.AddressDTO, error) {
	if err := model.DeleteAddress(userID, id); err != nil {
		return nil, err
	}
	return s.List(userID)
}

func (s *AddressService) SetDefault(userID uint, id uint) ([]dto.AddressDTO, error) {
	_, err := model.FindAddressByID(userID, id)
	if err != nil {
		return nil, NewError(404, "地址不存在")
	}
	if err := model.DB.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		return nil, err
	}
	if err := model.DB.Model(&model.Address{}).Where("user_id = ? AND id = ?", userID, id).Update("is_default", true).Error; err != nil {
		return nil, err
	}
	return s.List(userID)
}

func addressDTO(addr *model.Address) dto.AddressDTO {
	return dto.AddressDTO{ID: addr.ID, Receiver: addr.Receiver, Phone: addr.Phone, Province: addr.Province, City: addr.City, District: addr.District, Detail: addr.Detail, IsDefault: addr.IsDefault}
}
