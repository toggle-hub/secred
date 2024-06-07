package middlewares

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/xsadia/secred/pkg/utils"
)

var ErrInvalidSignMethod = errors.New("invalid signing method")
var ErrInvalidToken = errors.New("invalid token")

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ah := c.Request().Header.Get(echo.HeaderAuthorization)
		if ah == "" {
			return next(c)
		}

		split := strings.Split(ah, "Bearer ")
		tokenString := split[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secretKey := []byte(utils.Or(os.Getenv("JWT_SECRET"), "secret-key"))

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("Client error",
					ErrInvalidSignMethod)
				return nil, ErrInvalidSignMethod
			}

			return secretKey, nil
		})

		if err != nil {
			log.Println("Client error",
				err)
			return next(c)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			sub, ok := claims["sub"].(string)
			if !ok {
				return next(c)
			}

			c.Set("user", sub)
			return next(c)
		}

		return next(c)
	}
}
