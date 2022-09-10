package utils

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/v1/dtos"
)

func GetRequestData(c echo.Context) (dtos.JwtPayload, error) {
	var userData dtos.JwtPayload
	_token := c.Get("user")
	if _token == nil {
		return userData, errors.New("request data not found")
	}
	token := _token.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	dataMar, _ := json.Marshal(claims["data"])
	err := json.Unmarshal(dataMar, &userData)
	if err != nil {
		log.Println("requester data Unmarshal err:", err.Error())
		return userData, err
	}
	return userData, nil
}

func IsSuperAdmin(c echo.Context) bool {
	requesterData, err := GetRequestData(c)
	log.Println("Utils Requester file! Is Super admin:", err == nil && requesterData.Role == string(enums.ROLE_SUPER_ADMIN))
	return err == nil && requesterData.Role == string(enums.ROLE_SUPER_ADMIN)
}
