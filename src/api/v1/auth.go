package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sajalmia381/store-api/src/api/common"
	"github.com/sajalmia381/store-api/src/config"
	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/v1/api"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"github.com/sajalmia381/store-api/src/v1/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type authApi struct {
	authService  service.AuthService
	jwtService   service.JwtService
	tokenService service.TokenService
}

func (a authApi) Login(c echo.Context) error {
	var payload dtos.LoginPayload
	if err := c.Bind(&payload); err != nil {
		return common.GenerateErrorResponse(c, err.Error(), "Failed to bind data")
	}
	if err := payload.Validate(); err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	user, err := a.authService.Login(payload)
	if err != nil {
		return common.GenerateErrorResponse(c, "[ERROR]: User is not matched!", err.Error())
	}
	if !user.Status {
		return common.GenerateErrorResponse(c, "[ERROR]: You were disabled!", "Please contact to admin for active you account!")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return common.GenerateErrorResponse(c, "[ERROR] Password is not matched!", "Password is wrong!")
	}
	jwtPayload := dtos.JwtPayload{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}

	tokenExpiry, err := strconv.ParseInt(config.RegularTokenLifetime, 10, 64)
	if err != nil {
		return common.GenerateErrorResponse(c, "[ERROR]: failed to read regular token lifetime from env!", err.Error(), &common.ResponseOption{
			HttpCode: http.StatusInternalServerError,
		})
	}
	jwtRes, err := a.jwtService.GenerateToken(jwtPayload, tokenExpiry)
	if err != nil {
		return common.GenerateErrorResponse(c, "[ERROR]: failed to generate tokens!", err.Error(), &common.ResponseOption{
			HttpCode: http.StatusInternalServerError,
		})
	}
	if user.Role == enums.ROLE_SUPER_ADMIN {
		_, err = a.authService.UpdateUserLoginTime(user.ID)
		if err != nil {
			log.Println("Failed to update user login time", err.Error())
		}
		token := model.Token{
			ID:     primitive.NewObjectID(),
			UserId: &user.ID,
			Type:   string(enums.REFRESH_TOKEN),
			Token:  jwtRes.RefreshToken,
		}
		if _, err := a.tokenService.Store(token); err != nil {
			log.Println("Failed to store super admin refresh token:", err.Error())
		}
	}
	return common.GenerateSuccessResponse(c, jwtRes, "Success! User logged in")
}

func (a authApi) Register(c echo.Context) error {
	var payload dtos.RegisterPayload
	if err := c.Bind(&payload); err != nil {
		return common.GenerateErrorResponse(c, err.Error(), "Failed to bind data")
	}
	if err := payload.Validate(); err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	user, err := a.authService.Register(payload)
	if err != nil {
		return common.GenerateErrorResponse(c, "[ERROR]: failed to register user!", err.Error())
	}
	jwtPayload := dtos.JwtPayload{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
	tokenExpiry, err := strconv.ParseInt(config.RegularTokenLifetime, 10, 64)
	if err != nil {
		return common.GenerateErrorResponse(c, "[ERROR]: failed to read regular token lifetime from env!", err.Error(), &common.ResponseOption{
			HttpCode: http.StatusInternalServerError,
		})
	}
	jwtRes, err := a.jwtService.GenerateToken(jwtPayload, tokenExpiry)
	if err != nil {
		return common.GenerateErrorResponse(c, "[ERROR]: failed to generate tokens!", err.Error(), &common.ResponseOption{
			HttpCode: http.StatusInternalServerError,
		})
	}
	return common.GenerateSuccessResponse(c, jwtRes, "Success! User registration successful")
}

func (a authApi) RefreshToken(c echo.Context) error {
	var payload dtos.RefreshTokenPayload
	if err := c.Bind(&payload); err != nil {
		return common.GenerateErrorResponse(c, err.Error(), "Failed to bind data")
	}
	if err := payload.Validate(); err != nil {
		return common.GenerateErrorResponse(c, nil, err.Error())
	}
	isValid := a.jwtService.VerifyToken(payload.RefreshToken)
	if !isValid {
		return common.GenerateErrorResponse(c, "[ERROR]: Token is expired!", "Please login again to get token!", &common.ResponseOption{
			HttpCode: http.StatusForbidden,
		})
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(payload.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtService.GetRegularTokenSecret(), nil
	})
	jsonBody, err := json.Marshal(claims["data"])
	if err != nil {
		log.Println(err)
	}
	var jwtPayload dtos.JwtPayload
	err = json.Unmarshal(jsonBody, &jwtPayload)
	if err != nil {
		log.Println("[ERROR] failed to Unmarshal", err.Error())
		return common.GenerateErrorResponse(c, "[ERROR]: failed to get date from token!", err.Error(), &common.ResponseOption{
			HttpCode: http.StatusInternalServerError,
		})
	}
	tokenExpiry, err := strconv.ParseInt(config.RegularTokenLifetime, 10, 64)
	if err != nil {
		return common.GenerateErrorResponse(c, "[ERROR]: failed to read regular token lifetime from env!", err.Error(), &common.ResponseOption{
			HttpCode: http.StatusInternalServerError,
		})
	}
	jwtRes, err := a.jwtService.GenerateToken(jwtPayload, tokenExpiry)
	if err != nil {
		return common.GenerateErrorResponse(c, "[ERROR]: failed to generate tokens!", err.Error(), &common.ResponseOption{
			HttpCode: http.StatusInternalServerError,
		})
	}
	return common.GenerateSuccessResponse(c, jwtRes, "Success! New tokens generated")
}

func NewAuthApi(authService service.AuthService, jwtService service.JwtService, tokenService service.TokenService) api.AuthApi {
	return &authApi{
		authService:  authService,
		jwtService:   jwtService,
		tokenService: tokenService,
	}
}
