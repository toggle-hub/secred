package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const JWTExpireTime = 60 * 60 * 1000 * 24

func CreateJWT(id string, expireAt time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "secred",
		"sub": id,
		"exp": time.Now().Add(expireAt * time.Millisecond).Unix(),
	})

	key := Or(os.Getenv("JWT_SECRET"), "your-secret-key")

	signedToken, err := token.SignedString([]byte(key))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
