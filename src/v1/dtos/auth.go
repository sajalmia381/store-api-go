package dtos

import (
	"errors"
	"strconv"
)

type LoginPayload struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (p LoginPayload) Validate() error {
	if p.Email == "" {
		return errors.New("email is required")
	}
	if p.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type RegisterPayload struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Number   *uint  `json:"number" bson:"number"`
}

func (p RegisterPayload) Validate() error {
	if p.Email == "" {
		return errors.New("email is required")
	}
	if p.Password == "" {
		return errors.New("password is required")
	}
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Number != nil {
		if p.Number != nil {
			if len(strconv.Itoa((int(*p.Number)))) < 9 || len(strconv.Itoa((int(*p.Number)))) > 11 {
				return errors.New("number must be GREATER than 9 digit or LESS than 11 digit")
			}
		}
	}
	return nil
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

func (p RefreshTokenPayload) Validate() error {
	if p.RefreshToken == "" {
		return errors.New("refresh token is required")
	}
	return nil
}
