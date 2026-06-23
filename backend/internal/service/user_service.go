package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gomall-lite-api/config"
	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/model"
	"gorm.io/gorm"
)

type UserService struct {
	cfg config.Config
}

func NewUserService(cfg config.Config) *UserService {
	return &UserService{cfg: cfg}
}

func (s *UserService) Register(req dto.RegisterRequest) (*dto.UserDTO, error) {
	if req.Username == "" || req.Password == "" {
		return nil, NewError(400, "用户名和密码不能为空")
	}
	_, err := model.FindUserByUsername(req.Username)
	if err == nil {
		return nil, NewError(400, "用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if req.Nickname == "" {
		req.Nickname = req.Username
	}
	user := model.User{Username: req.Username, PasswordHash: string(hash), Nickname: req.Nickname}
	if err := model.CreateUser(&user); err != nil {
		return nil, err
	}
	return userDTO(&user), nil
}

func (s *UserService) Login(req dto.LoginRequest) (*dto.LoginDTO, error) {
	user, err := model.FindUserByUsername(req.Username)
	if err != nil {
		return nil, NewError(400, "用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, NewError(400, "用户名或密码错误")
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &dto.LoginDTO{Token: token, User: *userDTO(user)}, nil
}

func (s *UserService) GetUserInfo(userID uint) (*dto.UserDTO, error) {
	user, err := model.FindUserByID(userID)
	if err != nil {
		return nil, NewError(404, "用户不存在")
	}
	return userDTO(user), nil
}

func (s *UserService) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(s.cfg.TokenExpireHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func userDTO(user *model.User) *dto.UserDTO {
	return &dto.UserDTO{ID: user.ID, Username: user.Username, Nickname: user.Nickname}
}
