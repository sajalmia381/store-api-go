package dtos

import (
	"errors"
	"strconv"
	"time"
)

type (
	UserRegisterDTO struct {
		Name     string `json:"name" bson:"name"`
		Email    string `json:"email" bson:"email"`
		Password string `json:"password" bson:"password"`
		Number   *uint  `json:"number" bson:"number"`
	}
)

func (u *UserRegisterDTO) Validate() error {
	if u.Number != nil {
		if len(strconv.Itoa((int(*u.Number)))) < 9 || len(strconv.Itoa((int(*u.Number)))) > 11 {
			return errors.New("number must be GREATER than 9 digit or LESS than 11 digit")
		}
	}
	return nil
}

type UserUpdateDto struct {
	Name     string    `json:"name" bson:"name"`
	Password string    `json:"password" bson:"password"`
	Number   *uint     `json:"number" bson:"number"`
	Status   *bool     `json:"status" bson:"status"`
	UpdateAt time.Time `json:"updateAt" bson:"updateAt"`
}

func (u *UserUpdateDto) Validate() error {
	if u.Number != nil {
		if len(strconv.Itoa((int(*u.Number)))) < 9 || len(strconv.Itoa((int(*u.Number)))) > 11 {
			return errors.New("number must be GREATER than 9 digit or LESS than 11 digit")
		}
	}
	return nil
}

type UserQuery struct {
	Status *bool  `json:"status" bson:"status"`
	Role   string `json:"role" bson:"role"`
}
