package service

import (
	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/logger"
	"gomall-lite-api/internal/model"
)

type AddressService struct{}

func NewAddressService() *AddressService { return &AddressService{} }

func (s *AddressService) List(userID uint) ([]dto.AddressDTO, error) {
	addresses, err := model.ListAddresses(userID)
	if err != nil {
		logger.Default().Error("list addresses failed", "user_id", userID, "error", err)
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
			logger.Default().Error("create address reset default failed", "user_id", userID, "error", err)
			return nil, err
		}
	}
	addr := model.Address{UserID: userID, Receiver: req.Receiver, Phone: req.Phone, Province: req.Province, City: req.City, District: req.District, Detail: req.Detail, IsDefault: req.IsDefault}
	if err := model.CreateAddress(&addr); err != nil {
		logger.Default().Error("create address failed", "user_id", userID, "error", err)
		return nil, err
	}
	logger.Default().Info("create address success", "user_id", userID, "address_id", addr.ID, "is_default", addr.IsDefault)
	return s.List(userID)
}

func (s *AddressService) Update(userID uint, id uint, req dto.AddressRequest) ([]dto.AddressDTO, error) {
	addr, err := model.FindAddressByID(userID, id)
	if err != nil {
		logger.Default().Warn("update address failed: address not found", "user_id", userID, "address_id", id)
		return nil, NewError(404, "地址不存在")
	}
	if req.IsDefault {
		if err := model.DB.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
			logger.Default().Error("update address reset default failed", "user_id", userID, "address_id", id, "error", err)
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
		logger.Default().Error("update address save failed", "user_id", userID, "address_id", id, "error", err)
		return nil, err
	}
	logger.Default().Info("update address success", "user_id", userID, "address_id", id, "is_default", addr.IsDefault)
	return s.List(userID)
}

func (s *AddressService) Delete(userID uint, id uint) ([]dto.AddressDTO, error) {
	if err := model.DeleteAddress(userID, id); err != nil {
		logger.Default().Error("delete address failed", "user_id", userID, "address_id", id, "error", err)
		return nil, err
	}
	logger.Default().Info("delete address success", "user_id", userID, "address_id", id)
	return s.List(userID)
}

func (s *AddressService) SetDefault(userID uint, id uint) ([]dto.AddressDTO, error) {
	_, err := model.FindAddressByID(userID, id)
	if err != nil {
		logger.Default().Warn("set default address failed: address not found", "user_id", userID, "address_id", id)
		return nil, NewError(404, "地址不存在")
	}
	if err := model.DB.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		logger.Default().Error("set default address reset failed", "user_id", userID, "address_id", id, "error", err)
		return nil, err
	}
	if err := model.DB.Model(&model.Address{}).Where("user_id = ? AND id = ?", userID, id).Update("is_default", true).Error; err != nil {
		logger.Default().Error("set default address update failed", "user_id", userID, "address_id", id, "error", err)
		return nil, err
	}
	logger.Default().Info("set default address success", "user_id", userID, "address_id", id)
	return s.List(userID)
}

func addressDTO(addr *model.Address) dto.AddressDTO {
	return dto.AddressDTO{ID: addr.ID, Receiver: addr.Receiver, Phone: addr.Phone, Province: addr.Province, City: addr.City, District: addr.District, Detail: addr.Detail, IsDefault: addr.IsDefault}
}
