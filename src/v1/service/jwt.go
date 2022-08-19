package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sajalmia381/store-api/src/config"
	"github.com/sajalmia381/store-api/src/v1/dtos"
)

type JwtService interface {
	GenerateToken(payload dtos.JwtPayload, duration int64) (dtos.JwtResponseDto, error)
	VerifyToken(token string) bool
	GetRegularTokenSecret() []byte
	GetRefreshTokenSecret() []byte
}

type jwtService struct{}

func (s jwtService) GenerateToken(payload dtos.JwtPayload, duration int64) (dtos.JwtResponseDto, error) {
	var tokens dtos.JwtResponseDto

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": payload,
		"exp":  time.Now().UTC().Add(time.Duration(duration) * time.Millisecond).Unix(),
		"iat":  time.Now().UTC().Unix(),
	})
	// log.Println("Regular token secret", s.GetRegularTokenSecret())

	tokenStr, err := token.SignedString(s.GetRegularTokenSecret())
	if err != nil {
		log.Println("[ERROR] regular token generate:", err.Error())
		return tokens, err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": payload,
		"exp":  time.Now().UTC().Add(time.Duration(duration+duration/4) * time.Millisecond).Unix(),
		"iat":  time.Now().UTC().Unix(),
	})
	refreshTokenStr, err := refreshToken.SignedString(s.GetRegularTokenSecret())
	if err != nil {
		log.Println("[ERROR] refresh token generate:", err.Error())
		return tokens, err
	}
	tokens.AccessToken = tokenStr
	tokens.RefreshToken = refreshTokenStr
	return tokens, nil
}

func (s jwtService) VerifyToken(tokenString string) bool {
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.GetRegularTokenSecret(), nil
	})

	var tm time.Time
	switch iat := claims["exp"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}
	return time.Now().UTC().After(tm)
}

func (s jwtService) GetRegularTokenSecret() []byte {
	return []byte(config.JwtRegularSecretKey)
}

func (s jwtService) GetRefreshTokenSecret() []byte {
	return []byte(config.JwtRefreshSecretKey)
}

// func (s jwtService) GetRegularTokenSecret() *rsa.PublicKey {
// 	block, _ := pem.Decode([]byte(config.JwtRegularSecretKey))
// 	publicKeyImported, err := x509.ParsePKCS1PublicKey(block.Bytes)
// 	if err != nil {
// 		log.Print("ERROR x509 parse public key:", err.Error())
// 		panic(err)
// 	}
// 	log.Println("public x509 imported key:", publicKeyImported)
// 	return publicKeyImported
// }

// func (s jwtService) GetRefreshTokenSecret() *rsa.PrivateKey {
// 	block, _ := pem.Decode([]byte(config.JwtRefreshSecretKey))
// 	publicKeyImported, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		log.Print("ERROR parse private key:", err.Error())
// 		panic(err)
// 	}
// 	return publicKeyImported
// }

func NewJwtService() JwtService {
	return &jwtService{}
}
