package service

import (
	"errors"
	"log"
	"time"

	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Store(payload dtos.UserRegisterDTO) (model.User, error)
	FindAll(filterData dtos.UserQuery) ([]model.User, error)
	FindById(id string) (model.User, error)
	UpdateById(id string, payload dtos.UserUpdateDto) (model.User, error)
	DeleteById(id string) error
	// Fake Action
	FakeStore(payload dtos.UserRegisterDTO) (model.User, error)
	FakeUpdateById(id string, payload dtos.UserUpdateDto) (model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func (s userService) Store(payload dtos.UserRegisterDTO) (model.User, error) {
	var user model.User
	_, err := s.repo.FindByEmail(payload.Email)
	if err == nil {
		return user, errors.New("user is exists! try with another email")
	}
	user = mergePayloadDataToUser(payload)
	user.Role = enums.ROLE_CUSTOMER
	user.Status = true
	// No dependent
	user.ID = primitive.NewObjectID()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[ERROR] convert string to hash password:", err.Error())
		return user, err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	newUser, err := s.repo.Store(user)
	return newUser, err
}

func (s userService) FindAll(query dtos.UserQuery) ([]model.User, error) {
	objects, err := s.repo.FindAll(query)
	return objects, err
}

func (s userService) FindById(id string) (model.User, error) {
	var user model.User
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, errors.New("invalid user id")
	}
	user, err = s.repo.FindById(_id)
	return user, err
}

func (s userService) UpdateById(id string, formData dtos.UserUpdateDto) (model.User, error) {
	var user model.User
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}
	payload := primitive.M{}
	if formData.Name == "" {
		payload["name"] = user.Name
	}

	if formData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(formData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("[ERROR] convert string to hash password:", err.Error())
			panic(err)
		}
		payload["password"] = string(hashedPassword)
	}
	if formData.Number == nil {
		payload["number"] = user.Number
	}
	if formData.Status == nil {
		payload["status"] = &user.Status
	}
	newUser, err := s.repo.UpdateById(_id, payload)
	return newUser, err
}

func (s userService) DeleteById(id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user id")
	}
	_, err = s.repo.DeleteById(_id)
	return err
}

// None Super Admin
func (s userService) FakeStore(payload dtos.UserRegisterDTO) (model.User, error) {
	var user model.User
	_, err := s.repo.FindByEmail(payload.Email)
	if err == nil {
		return user, errors.New("user is exists! try with another email")
	}
	user = mergePayloadDataToUser(payload)
	user.Role = enums.Role("ROLE_CUSTOMER")
	user.Status = true
	user.ID = primitive.NewObjectID()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[ERROR] convert string to hash password:", err.Error())
		return user, err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	return user, err
}

func (s userService) FakeUpdateById(id string, payload dtos.UserUpdateDto) (model.User, error) {
	user, err := s.FindById(id)
	if err != nil {
		return user, err
	}
	user.UpdatedAt = time.Now().UTC()
	if payload.Name != "" {
		user.Name = payload.Name
	}
	if payload.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("[ERROR] convert string to hash password:", err.Error())
			panic(err)
		}
		user.Password = string(hashedPassword)
	}
	if payload.Number != nil {
		user.Number = payload.Number
	}
	if payload.Status != nil {
		user.Status = *payload.Status
	}
	return user, nil
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
