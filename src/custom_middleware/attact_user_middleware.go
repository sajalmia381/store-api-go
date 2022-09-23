package custom_middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sajalmia381/store-api/src/config"
)

func AttachUserMiddlewareConfig() middleware.JWTConfig {
	signingKey := []byte(config.JwtRegularSecretKey)
	config := middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			accessToken := c.Request().Header.Get("Authorization")
			return accessToken == ""
		},
		TokenLookup: "header:Authorization", // Default
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return signingKey, nil
			}

			// claims are of type `jwt.MapClaims` when token is created with `jwt.Parse`
			token, err := jwt.Parse(auth, keyFunc)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}
			return token, nil
		},
		ErrorHandler: func(err error) error {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Token is missing or expired or invalid",
			}
		},
	}
	return config
}
