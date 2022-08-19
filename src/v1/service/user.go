package service

import (
	"errors"

	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	repository "github.com/sajalmia381/store-api/src/v1/repository"
)

type UserService interface {
	Store(payload dtos.UserRegisterDTO) (model.User, error)
	FindAll() ([]model.User, error)
	FindById(id string) (model.User, error)
	UpdateById(id string, payload dtos.UserUpdateDto) (model.User, error)
	DeleteById(id string) error
}

type userService struct {
	repo repository.UserRepository
}

func (s userService) Store(payload dtos.UserRegisterDTO) (model.User, error) {

	user, err := s.repo.FindByEmail(payload.Email)
	if err == nil {
		return user, errors.New("user is exists! try with another email")
	}
	user = mergePayloadDataToUser(payload)
	user.Role = enums.Role("ROLE_CUSTOMER")
	user.Status = true
	newUser, err := s.repo.Store(user)
	return newUser, err
}

func (s userService) FindAll() ([]model.User, error) {
	objects, err := s.repo.FindAll()
	return objects, err
}

func (s userService) FindById(id string) (model.User, error) {
	user, err := s.repo.FindById(id)
	return user, err
}

func (s userService) UpdateById(id string, payload dtos.UserUpdateDto) (model.User, error) {
	user, err := s.repo.UpdateById(id, payload)
	return user, err
}

func (s userService) DeleteById(id string) error {
	_, err := s.repo.DeleteById(id)
	return err
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		repo: userRepo,
	}
}

// Helper
func mergePayloadDataToUser(payload dtos.UserRegisterDTO) model.User {
	user := model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		Number:   payload.Number,
	}
	return user
}
