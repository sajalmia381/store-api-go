package service

import (
	"errors"

	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	Login(payload dtos.LoginPayload) (model.User, error)
	Register(payload dtos.RegisterPayload) (model.User, error)
	RefreshToken(payload dtos.RefreshTokenPayload)
	UpdateUserLoginTime(id primitive.ObjectID) (model.User, error)
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func (s authService) Login(payload dtos.LoginPayload) (model.User, error) {
	user, err := s.userRepo.FindByEmail(payload.Email)
	if err != nil {
		return user, errors.New("user is not found")
	}

	return user, nil
}

func (s authService) Register(payload dtos.RegisterPayload) (model.User, error) {
	user, err := s.userRepo.FindByEmail(payload.Email)
	if err == nil {
		return user, errors.New("user is exists! try with another email")
	}
	user = mergeRegisterDataToUser(payload)
	user.Role = enums.Role("ROLE_CUSTOMER")
	user.Status = true
	newUser, err := s.userRepo.Store(user)
	return newUser, err
}

func (s authService) RefreshToken(payload dtos.RefreshTokenPayload) {
	panic("not implemented") // TODO: Implement
}

func (s authService) UpdateUserLoginTime(id primitive.ObjectID) (model.User, error) {
	user, err := s.userRepo.UpdateLoginTime(id)
	return user, err
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

// Helper
func mergeRegisterDataToUser(payload dtos.RegisterPayload) model.User {
	user := model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		Number:   payload.Number,
	}
	return user
}
