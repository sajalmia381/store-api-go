package service

import (
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenService interface {
	Store(payload model.Token) (model.Token, error)
	FindByToken(token string) (model.Token, error)
	DeleteByToken(token string) (*mongo.DeleteResult, error)
}

type tokenService struct {
	repo repository.TokenRepository
}

func (s tokenService) Store(payload model.Token) (model.Token, error) {
	token, err := s.repo.Store(payload)
	return token, err
}

func (s tokenService) FindByToken(token string) (model.Token, error) {
	tokenObj, err := s.repo.FindByToken(token)
	return tokenObj, err
}

func (s tokenService) DeleteByToken(token string) (*mongo.DeleteResult, error) {
	result, err := s.repo.DeleteByToken(token)
	return result, err
}

func NewTokenService(repo repository.TokenRepository) TokenService {
	return &tokenService{
		repo: repo,
	}
}
