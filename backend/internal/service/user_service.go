package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gomall-lite-api/config"
	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/logger"
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
	logger.Default().Info("user register attempt", "username", req.Username)
	if req.Username == "" || req.Password == "" {
		logger.Default().Warn("user register failed: empty username or password", "username", req.Username)
		return nil, NewError(400, "用户名和密码不能为空")
	}

	_, err := model.FindUserByUsername(req.Username)
	if err == nil {
		logger.Default().Warn("user register failed: username exists", "username", req.Username)
		return nil, NewError(400, "用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Default().Error("user register db lookup failed", "username", req.Username, "error", err)
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Default().Error("user register password hash failed", "username", req.Username, "error", err)
		return nil, err
	}
	if req.Nickname == "" {
		req.Nickname = req.Username
	}

	user := model.User{Username: req.Username, PasswordHash: string(hash), Nickname: req.Nickname}
	if err := model.CreateUser(&user); err != nil {
		logger.Default().Error("user register create failed", "username", req.Username, "error", err)
		return nil, err
	}

	logger.Default().Info("user register success", "user_id", user.ID, "username", user.Username)
	return userDTO(&user), nil
}

func (s *UserService) Login(req dto.LoginRequest) (*dto.LoginDTO, error) {
	logger.Default().Info("user login attempt", "username", req.Username)
	user, err := model.FindUserByUsername(req.Username)
	if err != nil {
		logger.Default().Warn("user login failed: user not found", "username", req.Username)
		return nil, NewError(400, "用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.Default().Warn("user login failed: wrong password", "user_id", user.ID, "username", user.Username)
		return nil, NewError(400, "用户名或密码错误")
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		logger.Default().Error("generate token failed", "user_id", user.ID, "error", err)
		return nil, err
	}

	logger.Default().Info("user login success", "user_id", user.ID, "username", user.Username)
	return &dto.LoginDTO{Token: token, User: *userDTO(user)}, nil
}

func (s *UserService) GetUserInfo(userID uint) (*dto.UserDTO, error) {
	user, err := model.FindUserByID(userID)
	if err != nil {
		logger.Default().Warn("get user info failed: user not found", "user_id", userID)
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
