package dtos

import "errors"

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
	Number   *int   `json:"number" bson:"number"`
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
